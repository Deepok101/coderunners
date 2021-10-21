package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Deepok101/coderunners/pkg/services"
	utils "github.com/Deepok101/coderunners/utils/queue"
	"github.com/gorilla/mux"
)

func runCodeHandler(w http.ResponseWriter, r *http.Request) {
	var c utils.Code
	cQueue := services.GetServices().CodeQueue
	err := json.NewDecoder(r.Body).Decode(&c)
	c.Output = make(chan string)
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

func debugSetupHandler(w http.ResponseWriter, r *http.Request) {
	var c utils.Code
	debugger := services.GetServices().Debugger
	err := json.NewDecoder(r.Body).Decode(&c)

	if c.Language == "" || c.Content == "" {
		http.Error(w, "Wrong input", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.Output = make(chan string)
	debugger.Debug(c)

	out, err := json.Marshal(debugger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(out)

}

func debugStepIn(w http.ResponseWriter, r *http.Request) {
	debugger := services.GetServices().Debugger
	err := debugger.StepIn()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(debugger)

	w.Write(out)
}

func debugStepOut(w http.ResponseWriter, r *http.Request) {
	debugger := services.GetServices().Debugger
	err := debugger.StepOut()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(debugger)

	w.Write(out)
}

func debugStepOver(w http.ResponseWriter, r *http.Request) {
	debugger := services.GetServices().Debugger
	err := debugger.StepOver()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(debugger)

	w.Write(out)
}

func debugSetBreakpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	newBreakpointStr := params["lineNo"]
	newBreakpointInt, err := strconv.ParseInt(newBreakpointStr, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	debugger := services.GetServices().Debugger

	err = debugger.SetBreakpoint(int(newBreakpointInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(fmt.Sprintf("New breakpoint added; %d", debugger.Breakpoints)))

}
