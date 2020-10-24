package settings

import (
	"os"
)

const (
	artifactsDirectory      string = "ARTIFACTS_DIRECTORY"
	uniqueJobIdentifier     string = "UNIQUE_JOB_ID"
	controllerRunnerChannel string = "NATS_CONTROLLER_RUNNER_CHANNEL"
)

// GetArtifactsDirectory path to artifacts directory on runner
func GetArtifactsDirectory() string {
	return os.Getenv(artifactsDirectory)
}

// GetUniqueJobIdentifier unique id for this runner
func GetUniqueJobIdentifier() string {
	return os.Getenv(uniqueJobIdentifier)
}

// GetControllerRunnerChannelName name of channel to talk to controller
func GetControllerRunnerChannelName() string {
	return os.Getenv(controllerRunnerChannel)
}
