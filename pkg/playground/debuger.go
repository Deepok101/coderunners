package playground

import (
	"fmt"
	"os"
	"os/exec"

	utils "github.com/Deepok101/coderunners/utils/queue"
)

type Debugger interface {
	DebugMode(utils.Code) Debugger
	Next()
	Restart()
	SetBreakpoint()
}

type debugger struct {
	variables map[string]interface{}
	functions string
	execPath  string
}

func (d debugger) setDebugModePython() {
	cmd := exec.Command("bash", "-c", "python3 -m pdb coderunners/coderunners.py")
	cmd.Stdout = os.Stdout

	fmt.Println(cmd.Run())
}
