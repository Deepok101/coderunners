package playground

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	utils "github.com/Deepok101/coderunners/utils/queue"
)

type Debugger struct {
	running        bool
	LocalVariables map[string]interface{} `json:"localVars"`
	LineNo         int                    `json:"lineNo"`
	Function       string                 `json:"funcName"`
	Breakpoints    []int
}

func NewDebugger() *Debugger {
	return &Debugger{running: false}
}

func (d *Debugger) SetBreakpoint(breakpoint int) error {
	d.Breakpoints = append(d.Breakpoints, breakpoint)
	if d.running == true {
		print("aaaa")
		_, err := http.Get(fmt.Sprintf("http://localhost:8000/debug/set_breakpoint/%d", breakpoint))
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Debugger) Debug(code utils.Code) {
	switch code.Language {
	case "python":
		err := d.setupPythonDebug()
		if err != nil {
			print(err.Error())
			return
		}
	}
}

func (d *Debugger) setupPythonDebug() error {
	cwd, err := os.Getwd()
	pathToDebugger := filepath.FromSlash(cwd + "/static/python_debugger/debugger_api.py")
	cmd := exec.Command("python", pathToDebugger)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()

	if err != nil {
		return err
	}
	requestBody, _ := json.Marshal(map[string][]int{
		"breakpoints": d.Breakpoints,
	})
	resp, err := http.Post("http://localhost:8000/debug/setup", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	err = json.Unmarshal(body, d)

	if err != nil {
		return err
	}

	d.running = true
	return nil
}

func (d *Debugger) StepIn() error {
	return d.step("stepin")
}

func (d *Debugger) StepOut() error {
	return d.step("stepout")
}

func (d *Debugger) StepOver() error {
	return d.step("stepover")
}

func (d *Debugger) step(stepType string) error {
	resp, err := http.Get("http://localhost:8000/debug/" + stepType)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, d)

	if err != nil {
		return err
	}

	return nil
}
