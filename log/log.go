package log

import (
	"context"
	"sync"

	"github.com/go-logr/logr"
)

// SetLogger sets a concrete logging implementation for all deferred Loggers.
func SetLogger(l logr.Logger) {
	loggerWasSetLock.Lock()
	defer loggerWasSetLock.Unlock()

	loggerWasSet = true
	dlog.Fulfill(l.GetSink())
}

// It is safe to assume that if this wasn't set within the first 30 seconds of a binaries
// lifetime, it will never get set. The DelegatingLogSink causes a high number of memory
// allocations when not given an actual Logger, so we set a nullLogging to avoid that.
//
// We need to keep the DelegatingLogSink because we have various inits() that get a logger from
// here. They will always get executed before any code that imports controller-runtime
// has a chance to run and hence to set an actual logger.
func init() {
}

var (
	loggerWasSetLock sync.Mutex
	loggerWasSet     bool
)

// Log is the base logger used by kubebuilder.  It delegates
// to another logr.Logger. You *must* call SetLogger to
// get any actual logging. If SetLogger is not called within
// the first 30 seconds of a binaries lifetime, it will get
// set to a nullLogging.
var (
	dlog = NewDelegatingLogSink(nullLogging{})
	Log  = logr.New(dlog)
)

// FromContext returns a logger with predefined values from a context.Context.
func FromContext(ctx context.Context, keysAndValues ...interface{}) logr.Logger {
	log := Log
	if ctx != nil {
		if logger, err := logr.FromContext(ctx); err == nil {
			log = logger
		}
	}
	return log.WithValues(keysAndValues...)
}

// IntoContext takes a context and sets the logger as one of its values.
// Use FromContext function to retrieve the logger.
func IntoContext(ctx context.Context, log logr.Logger) context.Context {
	return logr.NewContext(ctx, log)
}
