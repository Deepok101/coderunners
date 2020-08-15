package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type code struct {
	language string
	content  string
}

func runCodeHandler(w http.ResponseWriter, r *http.Request) {
	var c code
	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	NewCodeQueue()
}

func InitializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/run", runCodeHandler).Methods("POST")
}
