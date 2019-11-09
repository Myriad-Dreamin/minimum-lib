package logger

import (
	"fmt"
	std_logger_interface "github.com/Myriad-Dreamin/minimum-lib/logger/std-logger-interface"
	"log"
	"os"
)

type WrappedStdLoggerInterface interface {
	Fatalln(...interface{})
	Println(...interface{})
}

type wrappedStdLogger struct {
	std_logger_interface.StdLoggerOutputInterface
}

func (g wrappedStdLogger) Fatalln(v ...interface{}) {
	_ = g.Output(3, fmt.Sprintln(v...))
	os.Exit(1)
}

func (g wrappedStdLogger) Println(v ...interface{}) {
	_ = g.Output(3, fmt.Sprintln(v...))
}

type stdLogger struct {
	srcLogger WrappedStdLoggerInterface
	keyvals   []interface{}
}

func NewStdLogger(options ...interface{}) Logger {
	var lg stdLogger
	for i := range options {
		switch option := options[i].(type) {
		case std_logger_interface.StdLoggerOutputInterface:
			lg.srcLogger = wrappedStdLogger{option}
		case WrappedStdLoggerInterface:
			lg.srcLogger = option
		default:
			lg.keyvals = append(lg.keyvals, option)
		}
	}
	if lg.srcLogger == nil {
		lg.srcLogger = log.New(os.Stdout, "", log.Llongfile|log.Ldate)
	}

	return lg
}

func toKeyVals(keyvals ...interface{}) (kvs []interface{}) {
	l := len(keyvals)
	kvs = make([]interface{}, (l+1)>>1)
	if (l & 1) != 0 {
		kvs[l>>1] = keyvals[l-1]
		for i := 1; i < l; i += 2 {
			kvs[i>>1] = fmt.Sprintf("%v=%v", keyvals[i ^ 1], keyvals[i])
		}
	} else {
		for i := 0; i < l; i += 2 {
			kvs[i>>1] = fmt.Sprintf("%v=%v", keyvals[i], keyvals[i ^ 1])
		}
	}
	return
}

func (s stdLogger) args(msg, level string, keyvals ...interface{}) []interface{} {
	return append(append(
	append(append(make([]interface{}, 0, len(s.keyvals) + ((len(keyvals) + 1) >> 1) + 2), msg),
		toKeyVals("level", level)...), s.keyvals...), toKeyVals(keyvals...)...)
}

func (s stdLogger) Fatal(msg string, keyvals ...interface{}) {
	s.srcLogger.Fatalln(s.args(msg, "Fatal", keyvals...)...)
}

func (s stdLogger) Error(msg string, keyvals ...interface{}) {
	s.srcLogger.Println(s.args(msg, "Error", keyvals...)...)
}

func (s stdLogger) Debug(msg string, keyvals ...interface{}) {
	s.srcLogger.Println(s.args(msg, "Debug", keyvals...)...)
}

func (s stdLogger) Warn(msg string, keyvals ...interface{}) {
	s.srcLogger.Println(s.args(msg, "Warn", keyvals...)...)
}

func (s stdLogger) Info(msg string, keyvals ...interface{}) {
	s.srcLogger.Println(s.args(msg, "Info", keyvals...)...)
}

func (s stdLogger) With(keyvals ...interface{}) Logger {
	return stdLogger{keyvals: append(s.keyvals, toKeyVals(keyvals)...)}
}
