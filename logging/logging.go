package logging

import (
	"fmt"
	"log"
)

const (
	ErrorLevel   string = "ERROR"
	WarningLevel string = "WARN"
	InfoLevel    string = "INFO"
	DebugLevel   string = "DEBUG"
	TraceLevel   string = "TRACE"
)

func init() {
	log.SetFlags(0) // Removes all logger prefixes
}

func formatLogLevel(logLevel string) string {
	return fmt.Sprintf("[%s]", logLevel)
}

func prefixLogArgs(prefix string, args ...interface{}) []interface{} {
	return append([]interface{}{prefix}, args...)
}

func Error(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(ErrorLevel), args...)...)
}

func Errorf(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(ErrorLevel, args...)...)
}

func Warning(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(WarningLevel), args...)...)
}

func Warningf(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(WarningLevel, args...)...)
}

func Info(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(InfoLevel), args...)...)
}

func Infof(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(InfoLevel, args...)...)
}

func Debug(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(DebugLevel), args...)...)
}

func Debugf(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(DebugLevel, args...)...)
}

func Trace(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(TraceLevel), args...)...)
}

func Tracef(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(TraceLevel, args...)...)
}

func getLoggingFuncByLevel(level string) func(format string, args ...interface{}) {
	fun, found := map[string]func(format string, args ...interface{}){
		ErrorLevel:   Errorf,
		WarningLevel: Warningf,
		InfoLevel:    Infof,
		DebugLevel:   Debugf,
		TraceLevel:   Tracef,
	}[level]
	if !found {
		fun = Tracef // But really, please don't get here?
	}
	return fun
}
