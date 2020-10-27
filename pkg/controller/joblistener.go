package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/andreasM009/cloudshipper-agent/pkg/controller/configuration"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/definition"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/processor"
	"github.com/andreasM009/cloudshipper-agent/pkg/controller/runner"

	"github.com/andreasM009/cloudshipper-agent/pkg/channel"
	natsforwarder "github.com/andreasM009/cloudshipper-agent/pkg/controller/events/nats"
	natsio "github.com/nats-io/nats.go"
)

// JobListener listens for jobs to process
type JobListener struct {
	NatsChannel *channel.NatsChannel
	done        chan int
}

// DeploymentJob job to process
type DeploymentJob struct {
	DeploymentName string            `json:"deploymentName"`
	ID             string            `json:"id"`
	DefinitionID   string            `json:"definitionId"`
	Yaml           string            `json:"yaml"`
	Parameters     map[string]string `json:"parameters"`
	LiveStreamName string            `json:"liveStreamName"`
}

// NewJobListener new instance
func NewJobListener(natsChannel *channel.NatsChannel) *JobListener {
	return &JobListener{
		NatsChannel: natsChannel,
		done:        make(chan int, 1),
	}
}

// StartListeningAsync listens for jobs to
func (listener *JobListener) StartListeningAsync(ctx context.Context) error {
	_, err := listener.NatsChannel.NatsConn.Subscribe(listener.NatsChannel.NatsPublishName, func(msg *natsio.Msg) {
		job := DeploymentJob{}

		if err := json.Unmarshal(msg.Data, &job); err != nil {
			// todo, notify subscribers that something is wrong
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
			// todo, notify subscribers that something is wrong
			log.Println(err)
			return
		}

		// create runtime deployment
		dep, err := definition.NewFromDefinition(def, job.DefinitionID, job.DeploymentName, job.ID, job.Parameters)
		if err != nil {
			// todo, notify subscribers that something is wrong
			log.Println(err)
			return
		}

		// load runtime
		runtime := configuration.NewRuntime()
		runtime.FromFlags()

		// channels to runner
		runnerChannel, err := channel.NewNatsChannel(job.ID, runtime.NatsServerConnectionStrings, fmt.Sprintf("%s-cntrl", job.ID))
		if err != nil {
			log.Panic(err)
		}

		defer runnerChannel.Close()

		// live stream
		liveStream, err := channel.NewNatsStreamingChannel(job.LiveStreamName, runtime.NatsServerConnectionStrings, fmt.Sprintf("%s-live-cntrl", job.LiveStreamName), runtime.NatsStreamingClusterID, fmt.Sprintf("%s-live-cntrl", job.LiveStreamName))
		if err != nil {
			log.Panic(err)
		}

		defer liveStream.Close()

		// forwarder
		forwarder := natsforwarder.NewNatsStreamEventForwarder(liveStream)

		// Command processor
		proc := processor.NewCommandProcessor(runnerChannel, dep, forwarder)

		err = proc.ProcessAsync()
		if err != nil {
			log.Panic(err)
		}

		if runtime.IsKubernetes {
			// start the RunnerPod if runner is hosted in Kubernetes
			k8sclient := runner.CreateKubeClient()

			if k8sclient == nil {
				log.Panic()
			}

			runnerCtx, cancelRunner := context.WithCancel(context.Background())
			defer cancelRunner()

			k8srunner, err := runner.CreateAndWatchRunnerPod(runnerCtx, k8sclient, job.ID)
			if err != nil {
				log.Panic(err)
			}

			select {
			case exitcode := <-proc.Done():
				cancelRunner()
				<-k8srunner.Error()
				fmt.Println(fmt.Sprintf("Runner finished with exitcode: %d", exitcode))
				return
			case err := <-k8srunner.Error():
				if err != nil {
					log.Panic(err)
				}
			}
		} else {
			exitcode := <-proc.Done()
			fmt.Println(fmt.Sprintf("Runner finished with exitcode: %d", exitcode))
		}
	})

	if err != nil {
		return err
	}

	return nil
}

// Done done listening channel
func (listener *JobListener) Done() <-chan int {
	return listener.done
}
