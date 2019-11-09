package logger

import (
	kitlog "github.com/go-kit/kit/log"
	kitlevel "github.com/go-kit/kit/log/level"
	"os"
)

const (
	msgKey    = "_msg" // "_" prefixed to avoid collisions
)

type kitLogger struct {
	srcLogger kitlog.Logger
}

func NewKitLogger(logger kitlog.Logger, options ...interface{}) Logger {
	var lg = new(kitLogger)
	lg.srcLogger = logger

	if len(options) != 0 {
		lg = lg.With(options...).(*kitLogger)
	}
	return lg
}

// Info logs a message at level Info.
func (l *kitLogger) Info(msg string, keyvals ...interface{}) {
	if err := kitlog.With(kitlevel.Info(l.srcLogger), msgKey, msg).Log(keyvals...); err != nil {
		_ = kitlog.With(kitlevel.Error(l.srcLogger), msgKey, msg).Log("err", err)
	}
}

// Debug logs a message at level Debug.
func (l *kitLogger) Debug(msg string, keyvals ...interface{}) {
	if err := kitlog.With(kitlevel.Debug(l.srcLogger), msgKey, msg).Log(keyvals...); err != nil {
		_ = kitlog.With(kitlevel.Error(l.srcLogger), msgKey, msg).Log("err", err)
	}
}

// Warn logs a message at level Debug.
func (l *kitLogger) Warn(msg string, keyvals ...interface{}) {
	if err := kitlog.With(kitlevel.Debug(l.srcLogger), msgKey, msg).Log(keyvals...); err != nil {
		_ = kitlog.With(kitlevel.Error(l.srcLogger), msgKey, msg).Log("err", err)
	}
}

// Error logs a message at level Error.
func (l *kitLogger) Error(msg string, keyvals ...interface{}) {
	lWithMsg := kitlog.With(kitlevel.Error(l.srcLogger), msgKey, msg)
	if err := lWithMsg.Log(keyvals...); err != nil {
		_ = lWithMsg.Log("err", err)
	}
}

// Fatal logs a message at level Error.
func (l *kitLogger) Fatal(msg string, keyvals ...interface{}) {
	lWithMsg := kitlog.With(kitlevel.Error(l.srcLogger), msgKey, msg)
	if err := lWithMsg.Log(append(keyvals)...); err != nil {
		_ = lWithMsg.Log("err", err)
	}
	os.Exit(1)
}

// With returns a new contextual logger with keyvals prepended to those passed
// to calls to Info, Debug or Error.
func (l *kitLogger) With(keyvals ...interface{}) Logger {
	return &kitLogger{kitlog.With(l.srcLogger, keyvals...)}
}
