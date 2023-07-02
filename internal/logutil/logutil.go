// ---

package logutil

import (
	"fmt"
	"log"
)

var g Logger = &defaultLogger{}

type Logger interface {
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

// SetLogger replaces the default logger
func SetLogger(l Logger) {
	g = l
}

type defaultLogger struct{}

func (defaultLogger) Infof(msg string, args ...interface{}) {
	_ = log.Output(2, fmt.Sprintf(msg, args...))
}

func (defaultLogger) Errorf(msg string, args ...interface{}) {
	_ = log.Output(2, fmt.Sprintf(msg, args...))
}

func (defaultLogger) Fatalf(msg string, args ...interface{}) {
	_ = log.Output(2, fmt.Sprintf(msg, args...))
}

func Infof(msg string, args ...interface{}) {
	g.Infof(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	g.Errorf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	g.Fatalf(msg, args...)
}
