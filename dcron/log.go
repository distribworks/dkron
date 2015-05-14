package dcron

import (
	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Level = logrus.DebugLevel
}
