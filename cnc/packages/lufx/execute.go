package lufx

import (
	"bytes"
	"fmt"
	"github.com/valyala/fasttemplate"
	"io"
	"strings"
)

// ExecuteToString will print the output into a new string.
func (e *Executor) ExecuteToString(input string) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := e.Execute(input, buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Execute will print the output directly into the writer.
func (e *Executor) Execute(input string, writer io.Writer) error {
	input = strings.ReplaceAll(input, "\\x1b", "\x1b")

	t, err := fasttemplate.NewTemplate(input, e.startTag, e.endTag)
	if err != nil {
		return err
	}

	_, err = t.ExecuteFunc(writer, func(w io.Writer, tag string) (int, error) {
		var cmdArgs string
		el := strings.Split(tag, "(")
		for _, element := range el {
			cmdArgs = strings.Split(element, ")")[0]
		}

		cmdArgs = strings.ReplaceAll(cmdArgs, "\"", "")
		var args = strings.Split(strings.TrimSpace(cmdArgs), ",")

		if len(el) > 1 {
			el[0] += "()"
		}

		command, ok := e.commands[el[0]]
		if ok == false {
			return w.Write([]byte(fmt.Sprintf("unknown tag [%s]", el[0])))
		}

		return command(w, args)
	})

	return err

}
