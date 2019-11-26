package logger

import (
	"context"

	"cloud.google.com/go/logging"
)

//go:generate mockgen -destination=./mock/mock.go github.com/mkuchenbecker/brewery3/data/logger Log,LoggerGetter,Logger

type Scope string

type Log interface {
	Log(context.Context, string)
	Level(Severity) Log
	WithError(error) Log
	With(string, interface{}) Log
	LogIfError(ctx context.Context, err error, message string)
}

type Logger interface {
	Printf(string, ...interface{})
}

type LoggerGetter interface {
	Get(Severity) Logger
}

// Severity is the severity of the event described in a log entry. These
// guideline severity levels are ordered, with numerically smaller levels
// treated as less severe than numerically larger levels.
type Severity int

const (
	// Debug means debug or trace information.
	Debug = Severity(logging.Debug)
	// Info means routine information, such as ongoing status or performance.
	Info = Severity(logging.Info)
	// Warning means events that might cause problems.
	Warning = Severity(logging.Warning)
	// Error means events that are likely to cause problems.
	Error = Severity(logging.Error)
	// Critical means events that cause more severe problems or brief outages.
	Critical = Severity(logging.Critical)
)

func (s Severity) ToGoogleSeverity() logging.Severity {
	return logging.Severity(int(s))
}
