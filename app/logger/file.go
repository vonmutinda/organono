package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type fileLogger struct {
	logger *zap.Logger
}

func newFileLogger(filePath string) *fileLogger {

	if filePath == "" {
		return &fileLogger{}
	}

	logger := newZapLogger(filePath)

	return &fileLogger{
		logger: logger,
	}
}

func (l *fileLogger) Panic(msg string) {
	l.logger.Panic(msg)
}

func (l *fileLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panic(fmt.Sprintf(format, args...))
}

func (l *fileLogger) Fatal(msg string) {
	l.logger.Fatal(msg)
}

func (l *fileLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *fileLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *fileLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *fileLogger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *fileLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *fileLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *fileLogger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *fileLogger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *fileLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *fileLogger) Flush() error {
	if l.logger != nil {
		return l.logger.Sync()
	}

	return nil
}

func (l *fileLogger) Valid() bool {
	return l.logger != nil
}
