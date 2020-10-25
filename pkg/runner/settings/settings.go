package settings

import (
	"os"
)

const (
	artifactsDirectory string = "ARTIFACTS_DIRECTORY"
)

// GetArtifactsDirectory path to artifacts directory on runner
func GetArtifactsDirectory() string {
	return os.Getenv(artifactsDirectory)
}
