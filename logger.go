package lua4go

import (
	"io"

	"github.com/lunny/log"
)

type Logger interface {
	Infof(format string, content ...interface{})
	Info(content ...interface{})

	Errorf(format string, content ...interface{})
	Error(content ...interface{})

	Fatalf(format string, content ...interface{})
	Fatal(content ...interface{})
}

func NewLogger(out io.Writer) Logger {
	l := log.New(out, "[lua] ", log.Ldefault())
	l.SetOutputLevel(log.Ldebug)
	return l
}
