package main

import (
	"fmt"
	"net/http"

	"github.com/Deepok101/coderunners/pkg/queue"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	var codeQueue queue.Queue
	codeQueue = queue.NewCodeQueue()
	
}
