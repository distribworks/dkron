package logging

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type LogSplitter struct{}

// Levels returns levels that are supported by the hook
func (l *LogSplitter) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel}
}

// Fire handles logging events based on log levels.
func (l *LogSplitter) Fire(entry *logrus.Entry) error {
	if entry == nil {
		return fmt.Errorf("logrus entry is nil")
	}

	switch entry.Level {
	case logrus.WarnLevel, logrus.DebugLevel, logrus.InfoLevel, logrus.TraceLevel:
		entry.Logger.Out = os.Stdout
	case logrus.ErrorLevel, logrus.PanicLevel, logrus.FatalLevel:
		entry.Logger.Out = os.Stderr
	default:
		entry.Logger.Out = os.Stdout
	}
	return nil
}
