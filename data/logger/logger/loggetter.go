package logger

import (
	"cloud.google.com/go/logging"
	logger "github.com/mkuchenbecker/brewery3/data/logger"
)

type LogGetter struct {
	*logging.Logger
}

func (g *LogGetter) Get(sev logger.Severity) logger.Logger {
	return g.StandardLogger(sev.ToGoogleSeverity())
}
