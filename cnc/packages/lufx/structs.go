package lufx

import (
	"io"
	"sync"
)

type FxFunction func(session io.Writer, args []string) (int, error)

type Executor struct {
	commands map[string]FxFunction
	mutex    sync.Mutex

	startTag string
	endTag   string
}
