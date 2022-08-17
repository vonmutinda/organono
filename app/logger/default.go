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
	logger.Panic(msg)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Fatal(msg string) {
	logger.Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Error(msg string) {
	logger.Error(msg)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Flush() error {
	logger.Flush()
	return nil
}
