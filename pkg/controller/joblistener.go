package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/andreasM009/cloudshipper-agent/pkg/controller/events"

	"github.com/andreasM009/cloudshipper-agent/pkg/controller/configuration"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/definition"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/processor"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/runner"

	natsforwarder "github.com/andreasM009/cloudshipper-agent/pkg/controller/events/nats"
	"github.com/andreasM009/nats-library/channel"
	snatsio "github.com/nats-io/stan.go"
)

// JobListener listens for jobs to process
type JobListener struct {
	NatsStreamingChannel *channel.NatsStreamingChannel
	done                 chan int
	eventForwarder       events.EventForwarder
}

// DeploymentJob job to process
type DeploymentJob struct {
	TenantID       string            `json:"tenantId"`
	DeploymentName string            `json:"deploymentName"`
	ID             string            `json:"id"`
	DefinitionID   string            `json:"definitionId"`
	Yaml           string            `json:"yaml"`
	Parameters     map[string]string `json:"parameters"`
	LiveStreamName string            `json:"liveStreamName"`
}

// NewJobListener new instance
func NewJobListener(natsChannel *channel.NatsStreamingChannel, eventForwarder events.EventForwarder) *JobListener {
	return &JobListener{
		NatsStreamingChannel: natsChannel,
		done:                 make(chan int, 1),
		eventForwarder:       eventForwarder,
	}
}

// StartListeningAsync listens for jobs to
func (listener *JobListener) StartListeningAsync(ctx context.Context) error {
	var wg sync.WaitGroup
	var mustStop bool = false
	var isJobRunning bool = false
	var mtx sync.Mutex

	_, err := listener.NatsStreamingChannel.SnatNativeConnection.QueueSubscribe(listener.NatsStreamingChannel.NatsPublishName, "agents", func(msg *snatsio.Msg) {
		// check if listener must stop or not
		mtx.Lock()

		if mustStop {
			mtx.Unlock()
			return
		}

		isJobRunning = true
		wg.Add(1)
		mtx.Unlock()

		defer func() {
			mtx.Lock()
			defer mtx.Unlock()
			isJobRunning = false
		}()

		// now we can run the job, and notify waiter when it is finished
		defer wg.Done()

		job := DeploymentJob{}

		defer msg.Ack()

		if err := json.Unmarshal(msg.Data, &job); err != nil {
			// here we can only log an error, no information available to notify subscribers :-(
			log.Println(err)
			return
		}

		// replace all parameters
		yaml := job.Yaml

		for k, v := range job.Parameters {
			yaml = strings.ReplaceAll(yaml, k, v)
		}

		// create the definition
		def, err := definition.NewFromYaml([]byte(yaml))
		if err != nil {
			log.Println(err)
			listener.forwardErrorInDeploymentEvent(&job)
			return
		}

		// create runtime deployment
		dep, err := definition.NewFromDefinition(def, job.TenantID, job.DefinitionID, job.DeploymentName, job.ID, job.Parameters)
		if err != nil {
			log.Println(err)
			listener.forwardErrorInDeploymentEvent(&job)
			return
		}

		// load runtime
		runtime := configuration.NewRuntime()

		// channels to runner
		runnerChannel, err := channel.NewNatsChannelFromPool(job.ID, runtime.NatsConnectionName)
		if err != nil {
			log.Println(err)
			listener.forwardErrorInDeploymentEvent(&job)
			return
		}

		defer runnerChannel.Close()

		// live stream channel
		liveStream, err := channel.NewNatsStreamingChannelFromPool(job.LiveStreamName, runtime.NatsStreamingClusterID, runtime.NatsClientID)
		if err != nil {
			log.Println(err)
			listener.forwardErrorInDeploymentEvent(&job)
			return
		}

		defer liveStream.Close()

		// forwarder
		forwarder := natsforwarder.NewNatsStreamEventForwarder(liveStream, listener.eventForwarder)

		// Command processor
		proc := processor.NewCommandProcessor(runnerChannel, dep, forwarder)

		err = proc.ProcessAsync()
		if err != nil {
			log.Println(err)
			listener.forwardErrorInDeploymentEvent(&job)
			return
		}

		if runtime.IsKubernetes {
			// start the RunnerPod if runner is hosted in Kubernetes
			k8sclient := runner.CreateKubeClient()

			if k8sclient == nil {
				log.Println("Unable to start runner pod")
				listener.forwardErrorInDeploymentEvent(&job)
				// Todo: cancel processor
				return
			}

			runnerCtx, cancelRunner := context.WithCancel(context.Background())
			defer cancelRunner()

			k8srunner, err := runner.CreateAndWatchRunnerPod(runnerCtx, k8sclient, job.ID)
			if err != nil {
				// Todo: cancel processor
				log.Println(fmt.Sprintf("Failed to start and watch runnerpod: %s", err))
			}

			select {
			case exitcode := <-proc.Done():
				cancelRunner()
				<-k8srunner.Error()
				fmt.Println(fmt.Sprintf("Runner finished with exitcode: %d", exitcode))
				return
			case err := <-k8srunner.Error():
				if err != nil {
					log.Println(fmt.Sprintf("Unexpected error in runner pod: %s", err))
				}
			}
		} else {
			exitcode := <-proc.Done()
			log.Println(fmt.Sprintf("Runner finished with exitcode: %d", exitcode))
		}
	}, snatsio.DurableName("agent"), snatsio.SetManualAckMode(), snatsio.MaxInflight(1), snatsio.AckWait(time.Hour*12))

	if err != nil {
		return err
	}

	// wait until controller must stop
	go func() {
		<-ctx.Done()
		mtx.Lock()
		mustStop = true
		if isJobRunning {
			mtx.Unlock()
			wg.Wait()
		} else {
			mtx.Unlock()
		}

		listener.done <- -1
	}()

	return nil
}

// Done done listening channel
func (listener *JobListener) Done() <-chan int {
	return listener.done
}

func (listener *JobListener) forwardErrorInDeploymentEvent(job *DeploymentJob) {
	evt := events.DeploymentEvent{
		Event: events.Event{
			DeploymentName: job.DeploymentName,
			DefinitionID:   job.DefinitionID,
			DeploymentID:   job.DefinitionID,
			EventName:      "deploymentEvent",
		},
		Started:  false,
		Finished: true,
		Exitcode: -1,
	}

	listener.eventForwarder.ForwardDeploymentEvent(&evt)
}
