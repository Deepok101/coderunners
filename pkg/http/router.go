package http

import (
	"encoding/json"
	"net/http"

	"github.com/Deepok101/coderunners/pkg/queue"
	"github.com/gorilla/mux"
)

type code struct {
	language string
	content  string
}

var cQueue queue.Queue

func runCodeHandler(w http.ResponseWriter, r *http.Request) {
	var c code
	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cQueue.Enqueue(c)
}

// InitializeRouter initializes the main http router of the application
func InitializeRouter(codeQueue *queue.CodeQueue) {
	cQueue = codeQueue
	router := mux.NewRouter()
	router.HandleFunc("/run", runCodeHandler).Methods("POST")
}
