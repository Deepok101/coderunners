package services

import (
	"github.com/Deepok101/coderunners/pkg/docker"
	"github.com/Deepok101/coderunners/pkg/playground"
	utils "github.com/Deepok101/coderunners/utils/queue"
)

type Services struct {
	CodeQueue     utils.Queue
	DockerWrapper docker.CoderunnerDockerWrapper
	Playground    playground.Playground
}

var serviceSingleton Services

// NewServices creates a new instance of all the coderunner services located in the pkg folder. This creates a singleton
func NewServices() Services {

	serviceSingleton = Services{
		CodeQueue:     utils.NewCodeQueue(),
		DockerWrapper: docker.CreateNewCoderunnerDockerWrapper(),
		Playground:    playground.NewPlayground(""),
	}
	return serviceSingleton
}

// GetServices returns the Service singleton
func GetServices() Services {
	return serviceSingleton
}
