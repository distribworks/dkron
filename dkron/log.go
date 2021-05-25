package dkron

import (
	"io/ioutil"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ginOnce is a wrapper around gin global var changes. This is a workaround
// against the lack of concurrency safety of these vars in the gin package.
var ginOnce sync.Once

// InitLogger creates the logger instance
func InitLogger(logLevel string, node string) *logrus.Entry {
	formattedLogger := logrus.New()
	formattedLogger.Formatter = &logrus.TextFormatter{FullTimestamp: true}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Error("Error parsing log level, using: info")
		level = logrus.InfoLevel
	}

	formattedLogger.Level = level
	log := logrus.NewEntry(formattedLogger).WithField("node", node)

	ginOnce.Do(func() {
		if level == logrus.DebugLevel {
			gin.DefaultWriter = log.Writer()
			gin.SetMode(gin.DebugMode)
		} else {
			gin.DefaultWriter = ioutil.Discard
			gin.SetMode(gin.ReleaseMode)
		}
	})

	return log
}
