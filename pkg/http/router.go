package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// InitializeRouter initializes the main http router of the application
func InitializeRouter() {
	portNum := os.Getenv("MAIN_PORT")

	router := mux.NewRouter()
	router.HandleFunc("/run", runCodeHandler).Methods("POST")
	router.HandleFunc("/debug", debugSetupHandler).Methods("POST")
	router.HandleFunc("/debug/stepin", debugStepIn).Methods("GET")
	router.HandleFunc("/debug/stepout", debugStepOut).Methods("GET")
	router.HandleFunc("/debug/stepover", debugStepOver).Methods("GET")
	router.HandleFunc("/debug/setbreakpoint/{lineNo}", debugSetBreakpoint).Methods("GET")
	http.ListenAndServe(fmt.Sprintf(":%s", portNum), router)
}
