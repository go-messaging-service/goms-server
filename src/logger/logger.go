package logger

import (
	"fmt"
	"runtime"
	"strconv"
)

var DebugMode = false

func Debug(message string) {
	if !DebugMode {
		return
	}

	fmt.Printf("[debug] %s: %s\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message)
}

func Info(message string) {
	fmt.Printf("[info]  %s: %s\n", getCallerName(), message)
}

func Error(message string) {
	fmt.Errorf("[ERROR] %s: %s\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message)
}

func getCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() //runtime.FuncForPC(pc).Name()
}

func getCallerLine() int {
	_, _, line, _ := runtime.Caller(2)
	return line
}
