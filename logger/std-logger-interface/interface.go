package std_logger_interface


type StdLoggerOutputInterface interface {
	Output(callDepth int, msg string) error
}

type StdLoggerFatalInterface interface {
	Fatal(args ...interface{})
}

type StdLoggerPrintInterface interface {
	Print(args ...interface{})
}

type StdLoggerPanicInterface interface {
	Panic(args ...interface{})
}

type StdLoggerManFatalInterface interface {
	StdLoggerFatalInterface
	Fatalf(format string, args ...interface{})
	Fatalln(v ...interface{})
}

type StdLoggerManPrintInterface interface {
	StdLoggerPrintInterface
	Printf(format string, args ...interface{})
	Println(v ...interface{})
}

type StdLoggerManPanicInterface interface {
	StdLoggerPanicInterface
	Panicf(format string, args ...interface{})
	Panicln(v ...interface{})
}

type StdLoggerInterface interface {
	StdLoggerManPrintInterface
	StdLoggerManFatalInterface
	StdLoggerManPanicInterface
}
