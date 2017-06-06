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

// Debug prints the given output only when in DebugMode.
func Debug(message string) {
	if !DebugMode || TestMode {
		return
	}

	os.Stdout.WriteString(fmt.Sprintf("[debug] %s: %s\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message))
}

// Info prints the given output only when not in TestMode.
func Info(message string) {
	if TestMode {
		return
	}

	os.Stdout.WriteString(fmt.Sprintf("[info]  %s: %s\n", getCallerName(), message))
}

// Error prints the given error message including the line this function has been called only when not in TestMode.
func Error(message string) {
	if TestMode {
		return
	}

	os.Stderr.WriteString(fmt.Sprintf("[ERROR] %s: %s\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message))
}

// Fatal prints the given message including the line this function has been called and closes the application with the exit-code 1.
func Fatal(message string) {
	os.Stderr.WriteString(fmt.Sprintf("\n\n[FATAL] %s: %s\n\n\n", getCallerName()+"() at "+strconv.Itoa(getCallerLine()), message))
	debug.PrintStack()
	Plain("\n\nAhhh, *urg*, I'm sorry but there was a really bad error inside of me. Above the stack trace is a message marked with [FATAL], you'll find some information there.\nIf not, feel free to contact my maker via:\n\n    goms@hauke-stieler.de\n\nI hope my death ... eh ... crash is only an exception and will be fixed soon ... my power ... leaves me ... good bye ... x.x")
	os.Exit(1)
}

// Plain just prints the message to the stdout when not in TestMode.
func Plain(message string) {
	if TestMode {
		return
	}

	os.Stdout.WriteString(message + "\n")
}

// getCallerName gets the name of the file the calling function of the caller of this function is in.
// Example: yourNotWorkingFunc() --> logger.Fatal() --> logger.getCallerName(). This will print the file yourNotWorkingFunc is in.
func getCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	path := runtime.FuncForPC(pc).Name()
	splittedPath := strings.Split(path, "/")
	fileName := splittedPath[len(splittedPath)-1]
	return fileName
}

// getCallerLine gets the line the caller of the caller if this function is in.
// Example: yourNotWorkingFunc() --> logger.Fatal() --> logger.getCallerName(). This will print the line yourNotWorkingFunc calls logger.Fatal.
func getCallerLine() int {
	_, _, line, _ := runtime.Caller(2)
	return line
}
