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
var TestMode = false

func Debug(message string) {
	if !DebugMode || TestMode {
		return
	}

	os.Stdout.WriteString(fmt.Sprintf("[debug] %s: %s\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message))
}

func Info(message string) {
	if TestMode {
		return
	}

	os.Stdout.WriteString(fmt.Sprintf("[info]  %s: %s\n", getCallerName(), message))
}

func Error(message string) {
	if TestMode {
		return
	}

	os.Stderr.WriteString(fmt.Sprintf("[ERROR] %s: %s\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message))
}

func Fatal(message string) {
	os.Stderr.WriteString(fmt.Sprintf("\n\n[FATAL] %s: %s\n\n\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message))
	debug.PrintStack()
	Plain("\n\nAhhh, *urg*, I'm sorry but there was a really bad error inside of me. Above the stack trace is a message marked with [FATAL], you'll find some information there.\nIf not, feel free to contact my maker via:\n\n    goms@hauke-stieler.de\n\nI hope my death ... eh ... crash is only an exception and will be fixed soon ... my power ... leaves me ... good bye ... x.x")
	os.Exit(1)
}

func Plain(message string) {
	if TestMode {
		return
	}

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
