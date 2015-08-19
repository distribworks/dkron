package dkron

import (
	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.TextFormatter{FullTimestamp: true}
}
