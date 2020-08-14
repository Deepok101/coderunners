package docker

import (
	"fmt"

	"github.com/docker/docker/client"
)

func ConnectDocker() {
	//ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	fmt.Println(cli)
}
