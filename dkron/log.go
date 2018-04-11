package dkron

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

var log = logrus.NewEntry(logrus.New())

// InitLogger creates the logger instance
func InitLogger(logLevel string, node string) {
	formattedLogger := logrus.New()
	formattedLogger.Formatter = &logrus.TextFormatter{FullTimestamp: true}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Error("Error parsing log level, using: info")
		level = logrus.InfoLevel
	}

	formattedLogger.Level = level
	log = logrus.NewEntry(formattedLogger).WithField("node", node)
	gin.DefaultWriter = log.Writer()
}
