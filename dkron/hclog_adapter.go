package dkron

import (
	"bytes"
	"io"
	golog "log"

	"github.com/hashicorp/go-hclog"
	"github.com/sirupsen/logrus"
)

// HCLogAdapter implements the hclog interface, and wraps it
// around a Logrus entry
type HCLogAdapter struct {
	Logger     logrus.FieldLogger
	LoggerName string
}

// Log Emit a message and key/value pairs at a provided log level
func (*HCLogAdapter) Log(level hclog.Level, msg string, args ...interface{}) {
}

// Trace HCLog has one more level than we do. As such, we will never
// set trace level.
func (*HCLogAdapter) Trace(_ string, _ ...interface{}) {}

// Debug logging level message
func (a *HCLogAdapter) Debug(msg string, args ...interface{}) {
	a.CreateEntry(args).Debug(msg)
}

// Info logging level message
func (a *HCLogAdapter) Info(msg string, args ...interface{}) {
	a.CreateEntry(args).Info(msg)
}

// Warn logging level message
func (a *HCLogAdapter) Warn(msg string, args ...interface{}) {
	a.CreateEntry(args).Warn(msg)
}

// Error logging level message
func (a *HCLogAdapter) Error(msg string, args ...interface{}) {
	a.CreateEntry(args).Error(msg)
}

// IsTrace check
func (a *HCLogAdapter) IsTrace() bool {
	return false
}

// IsDebug check
func (a *HCLogAdapter) IsDebug() bool {
	return a.shouldEmit(logrus.DebugLevel)
}

// IsInfo check
func (a *HCLogAdapter) IsInfo() bool {
	return a.shouldEmit(logrus.InfoLevel)
}

// IsWarn check
func (a *HCLogAdapter) IsWarn() bool {
	return a.shouldEmit(logrus.WarnLevel)
}

// IsError check
func (a *HCLogAdapter) IsError() bool {
	return a.shouldEmit(logrus.ErrorLevel)
}

// SetLevel noop
func (a *HCLogAdapter) SetLevel(hclog.Level) {
	// interface definition says it is ok for this to be a noop if
	// implementations don't need/want to support dynamic level changing, which
	// we don't currently.
}

// With returns a new instance with the specified options
func (a *HCLogAdapter) With(args ...interface{}) hclog.Logger {
	e := a.CreateEntry(args)
	return &HCLogAdapter{Logger: e}
}

// Name returns the Name of the logger
func (a *HCLogAdapter) Name() string {
	return a.LoggerName
}

// Named returns a named logger
func (a *HCLogAdapter) Named(name string) hclog.Logger {
	var newName bytes.Buffer
	if a.LoggerName != "" {
		newName.WriteString(a.Name())
		newName.WriteString(".")
	}
	newName.WriteString(name)

	return a.ResetNamed(newName.String())
}

// ResetNamed returns a new logger with the default name
func (a *HCLogAdapter) ResetNamed(name string) hclog.Logger {
	fields := []interface{}{"subsystem_name", name}
	e := a.CreateEntry(fields)
	return &HCLogAdapter{Logger: e, LoggerName: name}
}

// StandardWriter return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (a *HCLogAdapter) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return nil
}

// StandardLogger is meant to return a stldib Logger type which wraps around
// hclog. It does this by providing an io.Writer and instantiating a new
// Logger. It then tries to interpret the log level by parsing the message.
//
// Since we are not using `hclog` in a generic way, and I cannot find any
// calls to this method from go-plugin, we will poorly support this method.
// Rather than pull in all of hclog writer parsing logic, pass it a Logrus
// writer, and hardcode the level to INFO.
//
// Apologies to those who find themselves here.
func (a *HCLogAdapter) StandardLogger(opts *hclog.StandardLoggerOptions) *golog.Logger {
	entry := a.Logger.WithFields(logrus.Fields{})
	return golog.New(entry.WriterLevel(logrus.InfoLevel), "", 0)
}

func (a *HCLogAdapter) shouldEmit(level logrus.Level) bool {
	currentLevel := a.Logger.WithFields(logrus.Fields{}).Level
	return currentLevel >= level
}

// CreateEntry creates a new logrus entry
func (a *HCLogAdapter) CreateEntry(args []interface{}) *logrus.Entry {
	if len(args)%2 != 0 {
		args = append(args, "<unknown>")
	}

	fields := make(logrus.Fields)
	for i := 0; i < len(args); i = i + 2 {
		k := args[i].(string)
		v := args[i+1]
		fields[k] = v
	}

	return a.Logger.WithFields(fields)
}

// ImpliedArgs returns With key/value pairs
func (a *HCLogAdapter) ImpliedArgs() []interface{} {
	return nil
}
