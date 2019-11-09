package sugar

import (
	"fmt"
	std_logger_interface "github.com/Myriad-Dreamin/minimum-lib/logger/std-logger-interface"
	"runtime"
)

func PrintStackToString() string {
	var buf [1024 * 10]byte
	n := runtime.Stack(buf[:], false)
	return fmt.Sprintf("==> %s\n", string(buf[:n]))
}

func PrintStack() {
	fmt.Print(PrintStackToString())
}

func PrintStackToPrinter(printInterface std_logger_interface.StdLoggerPrintInterface) {
	printInterface.Print(PrintStackToString())
}
