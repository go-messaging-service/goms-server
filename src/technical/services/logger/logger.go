package logger

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
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

func Fatal(message string) {
	os.Stderr.WriteString(fmt.Sprintf("\n\n[FATAL] %s: %s\n\n\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message))
	debug.PrintStack()
	os.Exit(1)
}

func Plain(message string) {
	os.Stdout.WriteString(message + "\n")
}

func getCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	path := runtime.FuncForPC(pc).Name()
	splittedPath := strings.Split(path, "/")
	fileName := splittedPath[len(splittedPath)-1]
	return fileName
}

func getCallerLine() int {
	_, _, line, _ := runtime.Caller(2)
	return line
}
