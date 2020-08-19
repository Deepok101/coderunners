package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Deepok101/coderunners/pkg/services"

	utils "github.com/Deepok101/coderunners/utils/queue"
	"github.com/gorilla/mux"
)

func runCodeHandler(w http.ResponseWriter, r *http.Request) {
	var c utils.Code
	cQueue := services.GetServices().CodeQueue
	err := json.NewDecoder(r.Body).Decode(&c)
	c.Output = make(chan string)
	fmt.Println(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cQueue.Enqueue(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	codeOutput := <-c.Output

	w.Write([]byte(codeOutput))
}

// InitializeRouter initializes the main http router of the application
func InitializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/run", runCodeHandler).Methods("POST")
	http.ListenAndServe(":80", router)
}
