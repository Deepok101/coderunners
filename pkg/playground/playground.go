package playground

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	utils "github.com/Deepok101/coderunners/utils/queue"
)

type Playground interface {
	ExecuteCode(utils.Code) (string, error)
}

type playground struct {
	execPath string
}

func NewPlayground(customPath string) Playground {
	p := playground{}
	if customPath == "" {
		folderName := "coderunners"
		cwd, err := os.Getwd()

		if err != nil {
			return nil
		}
		p.execPath = path.Join(cwd, folderName)

	} else {
		p.execPath = customPath
	}
	_, err := os.Stat(p.execPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(p.execPath, 0755)
		if err != nil {
			return nil
		}
	}

	err = os.Chdir(p.execPath)
	if err != nil {
		return nil
	}
	return p
}

func (p playground) ExecuteCode(code utils.Code) (string, error) {
	switch code.Language {
	case "python":
		out, err := p.executeCodePython(code)
		if err != nil {
			return "", err
		}
		return out, nil
	}

	return "", errors.New("Language is not supported yet")
}

func (p playground) executeCodePython(code utils.Code) (string, error) {
	scriptName := "coderunners.py"
	err := os.Chdir(p.execPath)
	if err != nil {
		return "", err
	}

	f, err := os.Create(scriptName)
	if err != nil {
		return "", err
	}

	_, err = fmt.Fprint(f, code.Content)
	if err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	buf := new(strings.Builder)
	cmd := exec.Command("python", path.Join(p.execPath, scriptName))
	cmd.Stdout = buf
	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
	return buf.String(), nil
}
