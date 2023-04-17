package log

import (
	"github.com/go-logr/logr"
)

type nullLogging struct{}

var _ logr.LogSink = nullLogging{}

func (log nullLogging) Init(logr.RuntimeInfo) {
}

func (nullLogging) Info(_ int, _ string, _ ...interface{}) {
	// Do nothing.
}

func (nullLogging) Enabled(_ int) bool {
	return false
}

func (nullLogging) Error(_ error, _ string, _ ...interface{}) {
	// Do nothing.
}

func (log nullLogging) WithName(_ string) logr.LogSink {
	return log
}

func (log nullLogging) WithValues(_ ...interface{}) logr.LogSink {
	return log
}
