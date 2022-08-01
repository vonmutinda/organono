package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZapLogger(outputFilePath string) *zap.Logger {

	level := zapcore.InfoLevel
	if os.Getenv("ENVIRONMENT") == "development" {
		level = zapcore.DebugLevel
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.Sampling = nil
	cfg.EncoderConfig = encoderConfig()

	if outputFilePath != "" {
		cfg.OutputPaths = []string{
			outputFilePath,
		}
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Failed while initializing zap logger err = %v", err)
	}

	return logger
}

func encoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()

	return zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		NameKey:      cfg.NameKey,
		LineEnding:   cfg.LineEnding,
		EncodeLevel:  cfg.EncodeLevel,
		EncodeTime:   zapcore.RFC3339NanoTimeEncoder,
		EncodeCaller: cfg.EncodeCaller,
		EncodeName:   cfg.EncodeName,
	}
}
