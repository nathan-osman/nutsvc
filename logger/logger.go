package logger

import (
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	EventServiceStatus = 1
	EventServerStatus  = 2
)

// Logger provides a central source for recording log messages.
type Logger struct {
	*eventlog.Log
}

// New creates a new logger.
func New(name string) (*Logger, error) {
	e, err := eventlog.Open(name)
	if err != nil {
		return nil, err
	}
	return &Logger{e}, nil
}
