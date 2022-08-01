package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type stdOutLogger struct {
	logger *zap.Logger
}

func newStdOutLogger() *stdOutLogger {
	return &stdOutLogger{
		logger: newZapLogger(""),
	}
}

func (l *stdOutLogger) Panic(msg string) {
	l.logger.Panic(msg)
}

func (l *stdOutLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panic(fmt.Sprintf(format, args...))
}

func (l *stdOutLogger) Fatal(msg string) {
	l.logger.Fatal(msg)
}

func (l *stdOutLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *stdOutLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *stdOutLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *stdOutLogger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *stdOutLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *stdOutLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *stdOutLogger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *stdOutLogger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *stdOutLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *stdOutLogger) Flush() error {
	if l.logger != nil {
		return l.logger.Sync()
	}
	return nil
}

func (l *stdOutLogger) Valid() bool {
	return l.logger != nil
}
