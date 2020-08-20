package services

import (
	"github.com/Deepok101/coderunners/pkg/docker"
	"github.com/Deepok101/coderunners/pkg/playground"
	utils "github.com/Deepok101/coderunners/utils/queue"
)

// Services : All except Debugger are interfaces. Debugger is a struct because we want to access its fields.
type Services struct {
	CodeQueue     utils.Queue
	DockerWrapper docker.CoderunnerDockerWrapper
	Playground    playground.Playground
	Debugger      *playground.Debugger
}

var serviceSingleton Services

// NewServices creates a new instance of all the coderunner services located in the pkg folder. This creates a singleton
func NewServices() Services {

	serviceSingleton = Services{
		CodeQueue:     utils.NewCodeQueue(),
		DockerWrapper: docker.CreateNewCoderunnerDockerWrapper(),
		Debugger:      playground.NewDebugger(),
		Playground:    playground.NewPlayground(""),
	}
	return serviceSingleton
}

// GetServices returns the Service singleton
func GetServices() Services {
	return serviceSingleton
}
