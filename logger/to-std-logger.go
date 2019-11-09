package logger

import (
	"fmt"
	std_logger_interface "github.com/Myriad-Dreamin/minimum-lib/logger/std-logger-interface"
)

type toStdLogger struct {
	srcLogger Logger
}

func (l toStdLogger) Print(args ...interface{}) {
	l.srcLogger.Info(fmt.Sprint(args...))
}

func (l toStdLogger) Printf(format string, args ...interface{}) {
	l.srcLogger.Info(fmt.Sprintf(format, args...))
}

func (l toStdLogger) Println(v ...interface{}) {
	l.srcLogger.Info(fmt.Sprintln(v...))
}

func (l toStdLogger) Fatal(args ...interface{}) {
	l.srcLogger.Fatal(fmt.Sprint(args...))
}

func (l toStdLogger) Fatalf(format string, args ...interface{}) {
	l.srcLogger.Fatal(fmt.Sprintf(format, args...))
}

func (l toStdLogger) Fatalln(v ...interface{}) {
	l.srcLogger.Fatal(fmt.Sprintln(v...))
}

func (l toStdLogger) Panic(args ...interface{}) {
	l.srcLogger.Fatal(fmt.Sprint(args...))
}

func (l toStdLogger) Panicf(format string, args ...interface{}) {
	l.srcLogger.Fatal(fmt.Sprintf(format, args...))
}

func (l toStdLogger) Panicln(v ...interface{}) {
	l.srcLogger.Fatal(fmt.Sprintln(v...))
}

func StdLogger(logger Logger) std_logger_interface.StdLoggerInterface {
	return toStdLogger{logger}
}




