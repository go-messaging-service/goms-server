package logger

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

var DebugMode = false

func Debug(message string) {
	if !DebugMode {
		return
	}

	os.Stdout.WriteString(fmt.Sprintf("[debug] %s: %s\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message))
}

func Info(message string) {
	os.Stdout.WriteString(fmt.Sprintf("[info]  %s: %s\n", getCallerName(), message))
}

func Error(message string) {
	os.Stderr.WriteString(fmt.Sprintf("[ERROR] %s: %s\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message))
}

func getCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() //runtime.FuncForPC(pc).Name()
}

func getCallerLine() int {
	_, _, line, _ := runtime.Caller(2)
	return line
}
