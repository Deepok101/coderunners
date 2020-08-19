package main

import (
	rest "github.com/Deepok101/coderunners/pkg/http"
	"github.com/Deepok101/coderunners/pkg/services"
)

func main() {
	services.NewServices()
	go executeCodeQueue()
	rest.InitializeRouter()
}

func executeCodeQueue() {
	services := services.GetServices()
	cQueue := services.CodeQueue
	codePlayground := services.Playground

	for {
		if cQueue.Length() != 0 {
			code, err := cQueue.Dequeue()
			if err != nil {
				code.Output <- err.Error()
				break
			}
			out, err := codePlayground.ExecuteCode(code)
			if err != nil {
				code.Output <- err.Error()
				return
			}
			code.Output <- out
		}
	}
}
