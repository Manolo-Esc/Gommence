package logger

import (
	"fmt"
	"runtime"

	"go.uber.org/zap"
)

var logger *zap.Logger

func Assert(condition bool) {
	_, file, line, callerInfoOk := runtime.Caller(1) // Obtener información de la pila de ejecución
	if logger == nil {
		if callerInfoOk {
			fmt.Printf("Error calling Assert: Logger not initialized at %s:%d (condition %t)\n", file, line, condition)
		} else {
			fmt.Printf("Error calling Assert: Logger not initialized and caller information unknown(condition %t)\n", condition)
		}
		return
	}
	assert(condition, callerInfoOk, file, line, "", false)
}

func AssertMessage(condition bool, message string) {
	_, file, line, callerInfoOk := runtime.Caller(1) // Obtener información de la pila de ejecución
	if logger == nil {
		if callerInfoOk {
			fmt.Printf("Error calling AssertMessage (%s): Logger not initialized at %s:%d (condition %t)\n", message, file, line, condition)
		} else {
			fmt.Printf("Error calling AssertMessage (%s): Logger not initialized and caller information unknown(condition %t)\n", message, condition)
		}
		return
	}
	assert(condition, callerInfoOk, file, line, message, false)
}

func assert(condition bool, callerInfoOk bool, file string, line int, message string, throwPanic bool) {
	if !condition {
		var logMessage string
		if message == "" {
			if callerInfoOk {
				logMessage = fmt.Sprintf("ASSERTION FAILED at %s:%d", file, line)
			} else {
				logMessage = fmt.Sprintf("ASSERTION FAILED at unknown location")
			}
		} else {
			if callerInfoOk {
				logMessage = fmt.Sprintf("ASSERTION FAILED at %s:%d - %s", file, line, message)
			} else {
				logMessage = fmt.Sprintf("ASSERTION FAILED at unknown location - %s", message)
			}
		}
		logger.Error(logMessage)
		if throwPanic {
			panic(logMessage)
		}
	}
}
