package main

import (
	rest "github.com/Deepok101/coderunners/pkg/http"
	"github.com/Deepok101/coderunners/pkg/services"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
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
			// TODO: To send the output of the code execution to the code.output channel could be executed in the ExecuteCode method
			if err != nil {
				code.Output <- err.Error()
				return
			}
			code.Output <- out
		}
	}
}
