package termfx

import (
	"advanced-telnet-cnc/packages/lufx"
	"os"
)

func (termFx *TermFX) Execute(path string, newLine bool, elements map[string]any) error {
	file, err := os.ReadFile(Directory + path)
	if err != nil {
		return termFx.session.Print("")
	}

	executor, err := lufx.NewExecutor(elements, "<<", ">>")

	toString, err := executor.ExecuteToString(string(file))
	if err != nil {
		return err
	}

	if len(toString) < 1 {
		return termFx.session.Print()
	}

	if newLine {
		return termFx.session.Print(toString + "\r\n")
	}

	return termFx.session.Print(toString)
}

func (termFx *TermFX) ExecuteString(path string, elements map[string]any) (string, error) {
	file, err := os.ReadFile(Directory + path)
	if err != nil {
		return "", err
	}

	executor, err := lufx.NewExecutor(elements, "<<", ">>")

	toString, err := executor.ExecuteToString(string(file))
	if err != nil {
		return "", err
	}

	if len(toString) < 1 {
		return "", nil
	}

	return toString, err
}
