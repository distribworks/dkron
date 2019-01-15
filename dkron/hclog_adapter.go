package dkron

import (
	"bytes"
	golog "log"

	"github.com/hashicorp/go-hclog"
	"github.com/sirupsen/logrus"
)

// HCLogAdapter implements the hclog interface, and wraps it
// around a Logrus entry
type HCLogAdapter struct {
	Log  logrus.FieldLogger
	Name string
}

// HCLog has one more level than we do. As such, we will never
// set trace level.
func (*HCLogAdapter) Trace(_ string, _ ...interface{}) {
	return
}

func (a *HCLogAdapter) Debug(msg string, args ...interface{}) {
	a.CreateEntry(args).Debug(msg)
}

func (a *HCLogAdapter) Info(msg string, args ...interface{}) {
	a.CreateEntry(args).Info(msg)
}

func (a *HCLogAdapter) Warn(msg string, args ...interface{}) {
	a.CreateEntry(args).Warn(msg)
}

func (a *HCLogAdapter) Error(msg string, args ...interface{}) {
	a.CreateEntry(args).Error(msg)
}

func (a *HCLogAdapter) IsTrace() bool {
	return false
}

func (a *HCLogAdapter) IsDebug() bool {
	return a.shouldEmit(logrus.DebugLevel)
}

func (a *HCLogAdapter) IsInfo() bool {
	return a.shouldEmit(logrus.InfoLevel)
}

func (a *HCLogAdapter) IsWarn() bool {
	return a.shouldEmit(logrus.WarnLevel)
}

func (a *HCLogAdapter) IsError() bool {
	return a.shouldEmit(logrus.ErrorLevel)
}

func (a *HCLogAdapter) SetLevel(hclog.Level) {
	// interface definition says it is ok for this to be a noop if
	// implementations don't need/want to support dynamic level changing, which
	// we don't currently.
}

func (a *HCLogAdapter) With(args ...interface{}) hclog.Logger {
	e := a.CreateEntry(args)
	return &HCLogAdapter{Log: e}
}

func (a *HCLogAdapter) Named(name string) hclog.Logger {
	var newName bytes.Buffer
	if a.Name != "" {
		newName.WriteString(a.Name)
		newName.WriteString(".")
	}
	newName.WriteString(name)

	return a.ResetNamed(newName.String())
}

func (a *HCLogAdapter) ResetNamed(name string) hclog.Logger {
	fields := []interface{}{"subsystem_name", name}
	e := a.CreateEntry(fields)
	return &HCLogAdapter{Log: e, Name: name}
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
	entry := a.Log.WithFields(logrus.Fields{})
	return golog.New(entry.WriterLevel(logrus.InfoLevel), "", 0)
}

func (a *HCLogAdapter) shouldEmit(level logrus.Level) bool {
	currentLevel := a.Log.WithFields(logrus.Fields{}).Level
	if currentLevel >= level {
		return true
	}

	return false
}

func (a *HCLogAdapter) CreateEntry(args []interface{}) *logrus.Entry {
	if len(args)%2 != 0 {
		args = append(args, "<unknown>")
	}

	fields := make(logrus.Fields)
	for i := 0; i < len(args); i = i + 2 {
		k, ok := args[i].(string)
		if !ok {
		}
		v := args[i+1]
		fields[k] = v
	}

	return a.Log.WithFields(fields)
}
