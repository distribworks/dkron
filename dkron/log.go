package dkron

import (
	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Formatter = &logrus.TextFormatter{FullTimestamp: true}
}

func SetLogLevel(debug bool) {
	if debug {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = logrus.InfoLevel
	}
}
