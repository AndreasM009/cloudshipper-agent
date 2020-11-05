package processor

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/andreasM009/cloudshipper-agent/pkg/requests"

	"github.com/andreasM009/cloudshipper-agent/pkg/logs"

	"github.com/andreasM009/cloudshipper-agent/pkg/controller/events"

	"github.com/andreasM009/cloudshipper-agent/pkg/commands"

	"github.com/andreasM009/cloudshipper-agent/pkg/controller/definition"
	"github.com/andreasM009/nats-library/channel"

	natsio "github.com/nats-io/nats.go"
)

// CommandProcessor processes commands on runner
type CommandProcessor struct {
	channel           *channel.NatsChannel
	runtimeDefinition *definition.RuntimeDeployment
	finishedChannel   chan int
	currentJobIdx     int
	currentCommandIdx int
	eventForwarder    events.EventForwarder
}

// NewCommandProcessor new instance
func NewCommandProcessor(channel *channel.NatsChannel, runtimeDefinition *definition.RuntimeDeployment, forwarder events.EventForwarder) *CommandProcessor {
	return &CommandProcessor{
		channel:           channel,
		runtimeDefinition: runtimeDefinition,
		finishedChannel:   make(chan int, 1),
		currentJobIdx:     0,
		currentCommandIdx: -1,
		eventForwarder:    forwarder,
	}
}

// ProcessAsync executes
func (processor *CommandProcessor) ProcessAsync() error {
	if len(processor.runtimeDefinition.Jobs) == 0 {
		return errors.New("No jobns defined to process")
	}

	processor.eventForwarder.ForwardDeploymentEvent(processor.createDeploymentEvent(true, false, 0))

	// subscribe to controller and runner msg channel and wait for runner requests
	_, err := processor.channel.NatsNativeConn.Subscribe(processor.channel.NatsPublishName, func(msg *natsio.Msg) {
		requestType, runnerRequest, err := getRequestFromStream(msg.Data)
		if err != nil {
			processor.eventForwarder.ForwardDeploymentEvent(processor.createDeploymentEvent(false, true, 1))
			log.Panic(err)
		}

		switch requestType {
		case requests.ControllerLog:
			if rq, ok := runnerRequest.(*requests.ControllerLogRequest); ok {
				processor.eventForwarder.ForwardCommandEvent(processor.createCommandEvent(&rq.Log))
			}
		case requests.ControllerQueryCommand:
			job := processor.runtimeDefinition.Jobs[processor.currentJobIdx]
			processor.currentCommandIdx++

			if processor.currentCommandIdx >= len(job.Commands) {
				msg.Respond(nil)
				processor.eventForwarder.ForwardDeploymentEvent(processor.createDeploymentEvent(false, true, 0))
				processor.finishedChannel <- 0
				return
			}

			step := job.Commands[processor.currentCommandIdx]
			if err := replyNextCommand(step.CommandType, &step.Command, msg); err != nil {
				log.Panic(err)
			}

		case requests.ControllerReportCommandError:
			if rq, ok := runnerRequest.(*requests.ControllerReportErrorRequest); ok {
				processor.eventForwarder.ForwardDeploymentEvent(processor.createDeploymentEvent(false, true, rq.Exitcode))
				processor.finishedChannel <- rq.Exitcode
				return
			}
		}
	})

	return err
}

// Done done?
func (processor *CommandProcessor) Done() <-chan int {
	return processor.finishedChannel
}

func getRequestFromStream(data []byte) (requests.RequestType, interface{}, error) {
	requestCarrier := &requests.RequestCarrier{}

	err := json.Unmarshal(data, requestCarrier)
	if err != nil {
		log.Panic(err) // todo
	}

	if requestCarrier.CarrierForType == requests.ControllerLog {
		buffer, err := json.Marshal(requestCarrier.Data)
		if err != nil {
			log.Panic(err)
		}

		rq := &requests.ControllerLogRequest{}
		err = json.Unmarshal(buffer, rq)
		if err != nil {
			log.Panic(err)
		}

		return requests.ControllerLog, rq, nil
	} else if requestCarrier.CarrierForType == requests.ControllerQueryCommand {
		return requests.ControllerQueryCommand, nil, nil
	} else if requestCarrier.CarrierForType == requests.ControllerReportCommandError {
		buffer, err := json.Marshal(requestCarrier.Data)
		if err != nil {
			log.Panic(err)
		}

		errrq := &requests.ControllerReportErrorRequest{}
		err = json.Unmarshal(buffer, errrq)
		if err != nil {
			log.Panic(err)
		}

		return requests.ControllerReportCommandError, errrq, nil
	} else {
		return requests.UnknownRequest, nil, errors.New("Received unknown command from runner")
	}
}

func replyNextCommand(commandType commands.CommandType, rq interface{}, msg *natsio.Msg) error {
	carrier := commands.CommandCarrier{
		CarrierForType: commandType,
		Data:           rq,
	}

	json, err := json.Marshal(carrier)

	if err != nil {
		log.Panic(err)
	}

	msg.Respond(json)
	return nil
}

func (processor *CommandProcessor) createCommandEvent(l *logs.LogMessage) *events.CommandEvent {
	job := processor.runtimeDefinition.Jobs[processor.currentJobIdx]
	cmd := job.Commands[processor.currentCommandIdx]

	evt := &events.CommandEvent{
		Event: events.Event{
			DefinitionID:   processor.runtimeDefinition.DefinitionID,
			DeploymentID:   processor.runtimeDefinition.ID,
			DeploymentName: processor.runtimeDefinition.DeploymentName,
			EventName:      "commandEvent",
			TenantID:       processor.runtimeDefinition.TenantID,
		},
		JobName:            job.Name,
		JobDisplayName:     job.Displayname,
		CommandName:        cmd.Name,
		CommandDisplayName: cmd.Displayname,
		CommandIndex:       cmd.Index,
		Logs:               []events.Log{events.NewLogFromLogs(l)},
	}

	return evt
}

func (processor *CommandProcessor) createDeploymentEvent(started, finished bool, exitcode int) *events.DeploymentEvent {
	jobs := make([]events.DeploymentJob, len(processor.runtimeDefinition.Jobs))

	for k, j := range processor.runtimeDefinition.Jobs {
		commands := make([]events.DeploymentCommand, len(j.Commands))

		for i, c := range j.Commands {
			ce := events.DeploymentCommand{
				Name:        c.Name,
				Displayname: c.Displayname,
				Index:       c.Index,
			}

			commands[i] = ce
		}

		job := events.DeploymentJob{
			Name:        j.Name,
			Displayname: j.Displayname,
			Commands:    commands,
		}

		jobs[k] = job
	}

	return &events.DeploymentEvent{
		Event: events.Event{
			DefinitionID:   processor.runtimeDefinition.DefinitionID,
			DeploymentID:   processor.runtimeDefinition.ID,
			DeploymentName: processor.runtimeDefinition.DeploymentName,
			EventName:      "deploymentEvent",
			TenantID:       processor.runtimeDefinition.TenantID,
		},
		Jobs:     jobs,
		Started:  started,
		Finished: finished,
		Exitcode: exitcode,
	}
}
