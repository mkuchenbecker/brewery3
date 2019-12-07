package logger

import (
	"context"
	"fmt"
	"log"
	"sync"

	logger "github.com/mkuchenbecker/brewery3/data/logger"
)

type standardLogger struct {
	get logger.Getter
	sev logger.Severity

	withMux sync.RWMutex
	with    map[string]interface{}
}

func New(get logger.Getter) logger.Log {
	return &standardLogger{
		get:  get,
		sev:  logger.Info,
		with: make(map[string]interface{}),
	}
}

func (log *standardLogger) Log(ctx context.Context, msg string) {
	toLog := msg
	log.withMux.RLock()
	defer log.withMux.RUnlock()
	for k, v := range log.with {
		toLog = fmt.Sprintf("%s\n\t%s: %+v", toLog, k, v)
	}
	l := log.get.Get(log.sev)
	l.Printf(toLog)
}

func (log *standardLogger) Level(sev logger.Severity) logger.Log {
	log.sev = sev
	return log
}

func (log *standardLogger) WithError(err error) logger.Log {
	if err == nil {
		return log
	}
	return log.With("error", err)
}

func (log *standardLogger) With(key string, val interface{}) logger.Log {
	log.withMux.Lock()
	defer log.withMux.Unlock()
	log.with[key] = val
	return log
}

func (log *standardLogger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	log.get.Get(logger.Info).Printf(format, args...)
}

func (log *standardLogger) LogIfError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	log.Level(logger.Error).WithError(err).Log(ctx, message)
}

type stdLogger struct{}

func NewSTD() logger.Log {
	return New(&stdLogger{})
}

func (l *stdLogger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	log.Printf(format, args...)
}

func (l *stdLogger) Get(logger.Severity) logger.Logger {
	return l
}
