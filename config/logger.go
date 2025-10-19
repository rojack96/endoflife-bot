package config

import (
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
)

var log *zap.Logger

// InitLogger initializes the logger with custom format
func InitLogger() {
	config := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			MessageKey:     "message",
			EncodeTime:     customTimeEncoder,
			EncodeLevel:    customLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

// Custom time encoder for [TIMESTAMP] format
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}

// Custom level encoder for [LEVEL] format with colors
func customLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var levelStr string
	switch l {
	case zapcore.DebugLevel:
		levelStr = "[" + colorBlue + "DEBUG" + colorReset + "]"
	case zapcore.InfoLevel:
		levelStr = "[" + colorGreen + "INFO" + colorReset + "]"
	case zapcore.WarnLevel:
		levelStr = "[" + colorYellow + "WARN" + colorReset + "]"
	case zapcore.ErrorLevel:
		levelStr = "[" + colorRed + "ERROR" + colorReset + "]"
	case zapcore.FatalLevel:
		levelStr = "[" + colorPurple + "FATAL" + colorReset + "]"
	default:
		levelStr = "[" + strings.ToUpper(l.String()) + "]"
	}
	enc.AppendString(levelStr)
}

// GetLogger returns the configured logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		InitLogger()
	}
	return log
}

// Helper functions for different log levels
func Info(message string, fields ...zap.Field) {
	GetLogger().Info(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	GetLogger().Error(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	GetLogger().Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	GetLogger().Warn(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	GetLogger().Fatal(message, fields...)
}
