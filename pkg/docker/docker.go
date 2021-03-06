package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	utils "github.com/Deepok101/coderunners/utils/queue"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// CoderunnerDockerWrapper is a Docker client wrapper that manages all the work related to playing with docker containers
type CoderunnerDockerWrapper interface {
	createClient() error
	ExecuteCode(utils.Code) (string, error)
}

type coderunnerDockerWrapper struct {
	dockerClient *client.Client
}

// CreateNewCoderunnerDockerWrapper creates a new empty CoderunnerDockerWrappe instance
func CreateNewCoderunnerDockerWrapper() CoderunnerDockerWrapper {
	newDockerwrapper := &coderunnerDockerWrapper{}
	newDockerwrapper.createClient()
	return newDockerwrapper
}

func (c *coderunnerDockerWrapper) createClient() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	c.dockerClient = cli
	return nil
}

func (c *coderunnerDockerWrapper) ExecuteCode(code utils.Code) (string, error) {
	switch code.Language {
	case "python":
		out, err := c.executeCodePython(code)
		if err != nil {
			return "", err
		}
		return out, nil
	}

	return "", errors.New("Language is not supported yet")
}

func (c *coderunnerDockerWrapper) executeCodePython(code utils.Code) (string, error) {
	ctx := context.Background()
	reader, err := c.dockerClient.ImagePull(ctx, "python", types.ImagePullOptions{})

	io.Copy(os.Stdout, reader)

	listOfCommands := []string{
		"touch coderunners.py",
		fmt.Sprintf("echo \"%s\" >> coderunners.py", code.Content),
		"cat coderunners.py",
		"python coderunners.py",
	}

	commands := strings.Join(listOfCommands, ";")

	res, err := c.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "python",
		Cmd:   []string{"sh", "-c", commands},
		Tty:   false,
	}, nil, nil, nil, "")

	if err != nil {
		return "", err
	}

	if err := c.dockerClient.ContainerStart(ctx, res.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := c.dockerClient.ContainerWait(ctx, res.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	out, err := c.dockerClient.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}

	buffer := new(strings.Builder)
	errBuffer := new(strings.Builder)

	stdcopy.StdCopy(buffer, errBuffer, out)

	return buffer.String(), nil

}

func CreateAndRunDockerContainer() error {
	ctx := context.Background()
	imageName := "ubuntu"
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	io.Copy(os.Stdout, reader)

	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"echo", "hello world"},
		Tty:   false,
	}, nil, nil, nil, "")

	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, res.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, res.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}
