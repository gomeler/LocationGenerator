package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//New should in theory let me use a standardized logger across all logging needs. There might be a better way to do this, but this does seem to work for now.
func New() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)
	return logger
}

func SetLevelDebug(theLogger *log.Logger) {
	theLogger.SetLevel(log.DebugLevel)
}

func SetLevelInfo(theLogger *log.Logger) {
	theLogger.SetLevel(log.InfoLevel)
}
