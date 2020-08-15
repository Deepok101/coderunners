package main

import (
	"fmt"
	"log"
	"net/http"

	docker "github.com/Deepok101/coderunners/pkg/Docker"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {

	docker.ConnectDocker()
	router := mux.NewRouter()
	router.HandleFunc("/", handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
