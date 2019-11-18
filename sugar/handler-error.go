package sugar

import (
	std_logger_interface "github.com/Myriad-Dreamin/minimum-lib/logger/std-logger-interface"
)

func HandlerError0(err error) {
	if err != nil {
		PrintStack()
		panic(err)
	}
}

func HandlerError(i interface{}, err error) interface{} {
	HandlerError0(err)
	return i
}

func HandlerError2(i, i2 interface{}, err error) []interface{} {
	HandlerError0(err)
	return []interface{}{i, i2}
}

func HandlerError3(i, i2, i3 interface{}, err error) []interface{} {
	HandlerError0(err)
	return []interface{}{i, i2, i3}
}

type HandlerErrorLogger struct {
	logger std_logger_interface.StdLoggerFatalInterface
}

func NewHandlerErrorLogger(logger std_logger_interface.StdLoggerFatalInterface) HandlerErrorLogger {
	return HandlerErrorLogger{logger: logger}
}

func (h HandlerErrorLogger) HandlerError0(err error) {
	if err != nil {
		PrintStack()
		h.logger.Fatal(err)
	}
}

func (h HandlerErrorLogger) HandlerError(i interface{}, err error) interface{} {
	h.HandlerError0(err)
	return i
}

func (h HandlerErrorLogger) HandlerError2(i, i2 interface{}, err error) []interface{} {
	h.HandlerError0(err)
	return []interface{}{i, i2}
}

func (h HandlerErrorLogger) HandlerError3(i, i2, i3 interface{}, err error) []interface{} {
	h.HandlerError0(err)
	return []interface{}{i, i2, i3}
}

