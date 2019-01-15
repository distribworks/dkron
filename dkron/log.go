package dkron

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.NewEntry(logrus.New())

// InitLogger creates the logger instance
func InitLogger(logLevel string, node string) logrus.FieldLogger {
	formattedLogger := logrus.New()
	formattedLogger.Formatter = &logrus.TextFormatter{FullTimestamp: true}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Error("Error parsing log level, using: info")
		level = logrus.InfoLevel
	}

	formattedLogger.Level = level
	log = logrus.NewEntry(formattedLogger).WithField("node", node)

	if level == logrus.DebugLevel {
		gin.DefaultWriter = log.Writer()
		gin.SetMode(gin.DebugMode)
	} else {
		gin.DefaultWriter = ioutil.Discard
		gin.SetMode(gin.ReleaseMode)
	}

	return log
}
