package main

import (
	"fmt"
	"net/http"

	"github.com/Deepok101/coderunners/pkg/docker"
	utils "github.com/Deepok101/coderunners/utils/queue"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	// var cQueue utils.Queue
	// cQueue = utils.NewCodeQueue()
	crDocker := docker.CreateNewCoderunnerDockerWrapper()
	crDocker.CreateClient()
	crDocker.ExecuteCode(utils.Code{Language: "python", Content: "print('a')"})

}
