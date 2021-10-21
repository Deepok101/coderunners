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
	Running        bool                   `json: running`
	LocalVariables map[string]interface{} `json:"localVars"`
	LineNo         int                    `json:"lineNo"`
	Function       string                 `json:"funcName"`
	Breakpoints    []int                  `json: breakpoints`
	Code           string
	urlMap         map[string]string
}

func NewDebugger() *Debugger {
	urlMap := make(map[string]string)
	urlMap["python"] = fmt.Sprintf("%s:%s/debug", os.Getenv("PYTHON_DEBUGGER_HOST"), os.Getenv("PYTHON_DEBUGGER_PORT"))
	fmt.Print(urlMap["python"])
	return &Debugger{Running: false, Breakpoints: []int{}, urlMap: urlMap}
}

func (d *Debugger) SetBreakpoint(breakpoint int) error {
	d.Breakpoints = append(d.Breakpoints, breakpoint)
	if d.Running {
		reqUrl := fmt.Sprintf("%s/set_breakpoint/%d", d.urlMap["python"], breakpoint)
		_, err := http.Get(reqUrl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Debugger) Debug(code utils.Code) {
	switch code.Language {
	case "python":
		err := d.setupPythonDebug(code)
		if err != nil {
			return
		}
	}
}

func (d *Debugger) setupPythonDebug(code utils.Code) error {
	cwd, err := os.Getwd()
	// Write to template file
	templateFileContent := "def coderunners_exec():\n" + code.Content
	ioutil.WriteFile(cwd+"/debugger/python_debugger/template.py", []byte(templateFileContent), 0644)

	pathToDebugger := filepath.FromSlash(cwd + "/debugger/python_debugger/debugger_api.py")
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
	// Call Python Debugger API
	reqUrl := fmt.Sprintf("%s/setup", d.urlMap["python"])
	resp, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(requestBody))
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
	reqUrl := fmt.Sprintf("%s/%s", d.urlMap["python"], stepType)
	resp, err := http.Get(reqUrl)

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
