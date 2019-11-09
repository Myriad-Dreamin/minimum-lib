package logger

type Logger interface {
	Fatal(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})

	With(keyvals ...interface{}) Logger
}



