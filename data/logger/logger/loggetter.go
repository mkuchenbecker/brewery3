package logger

import (
	"cloud.google.com/go/logging"
	logger "github.com/mkuchenbecker/brewery3/data/logger"
)

type Getter struct {
	*logging.Logger
}

func (g *Getter) Get(sev logger.Severity) logger.Logger {
	return g.StandardLogger(sev.ToGoogleSeverity())
}
