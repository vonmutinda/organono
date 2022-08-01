package logger

import (
	"os"
	"sync"
)

var (
	logger     Logger
	loggerSync = &sync.Once{}
)

func init() {
	logger = Default()
}

func Initialize() {
	loggerSync = &sync.Once{}
	logger = Default()
}

func Default() Logger {
	return setupLogger()
}

func setupLogger() Logger {
	loggerSync.Do(func() {
		var defaultLogger Logger = newFileLogger(os.Getenv("LOG_FILE"))
		if !defaultLogger.Valid() {
			defaultLogger = newStdOutLogger()
		}
		logger = defaultLogger
	})

	return logger
}

func Panic(msg string) {
	Default().Panic(msg)
}

func Panicf(format string, args ...interface{}) {
	Default().Panicf(format, args...)
}

func Fatal(msg string) {
	Default().Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	Default().Fatalf(format, args...)
}

func Error(msg string) {
	Default().Error(msg)
}

func Errorf(format string, args ...interface{}) {
	Default().Errorf(format, args...)
}

func Warn(msg string) {
	Default().Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	Default().Warnf(format, args...)
}

func Info(msg string) {
	Default().Info(msg)
}

func Infof(format string, args ...interface{}) {
	Default().Infof(format, args...)
}

func Debug(msg string) {
	Default().Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	Default().Debugf(format, args...)
}

func Flush() error {
	Default().Flush()
	return nil
}
