package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// InitializeRouter initializes the main http router of the application
func InitializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/run", runCodeHandler).Methods("POST")
	router.HandleFunc("/debug", debugSetupHandler).Methods("POST")
	router.HandleFunc("/debug/stepin", debugStepIn).Methods("GET")
	router.HandleFunc("/debug/stepout", debugStepOut).Methods("GET")
	router.HandleFunc("/debug/stepover", debugStepOver).Methods("GET")
	router.HandleFunc("/debug/setbreakpoint/{lineNo}", debugSetBreakpoint).Methods("GET")
	http.ListenAndServe(":80", router)
}
