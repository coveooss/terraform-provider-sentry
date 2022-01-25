package logging

import (
	"fmt"
	"log"
)

type LogLevel int

func (level LogLevel) String() string {
	return []string{"ERROR", "WARN", "INFO", "DEBUG", "TRACE"}[level]
}

const (
	ErrorLevel LogLevel = iota
	WarningLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func init() {
	log.SetFlags(0) // Removes all logger prefixes
}

func formatLogLevel(logLevel LogLevel) string {
	return fmt.Sprintf("[%s]", logLevel)
}

func prefixLogArgs(logLevel string, args ...interface{}) []interface{} {
	return append([]interface{}{logLevel}, args...)
}

func Error(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(ErrorLevel), args...)...)
}

func Errorf(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(ErrorLevel.String(), args...)...)
}

func Warning(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(WarningLevel), args...)...)
}

func Warningf(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(WarningLevel.String(), args...)...)
}

func Info(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(InfoLevel), args...)...)
}

func Infof(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(InfoLevel.String(), args...)...)
}

func Debug(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(DebugLevel), args...)...)
}

func Debugf(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(DebugLevel.String(), args...)...)
}

func Trace(args ...interface{}) {
	log.Println(prefixLogArgs(formatLogLevel(TraceLevel), args...)...)
}

func Tracef(format string, args ...interface{}) {
	log.Printf("[%s] "+format, prefixLogArgs(TraceLevel.String(), args...)...)
}

func getLoggingFuncByLevel(level LogLevel) func(format string, args ...interface{}) {
	return []func(format string, args ...interface{}){
		Errorf,
		Warningf,
		Infof,
		Debugf,
		Tracef,
	}[level]
}
