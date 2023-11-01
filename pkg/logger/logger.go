package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	Logger *zap.Logger
}

func New() *Logger {
	logger, err := zap.NewDevelopment(zap.WithCaller(true), zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic(err)
	}
	return &Logger{Logger: logger}
}

// Debug -.
func (l *Logger) Debug(message string, fields ...zap.Field) {
	l.Logger.Debug(message, fields...)
}

// Info -.
func (l *Logger) Info(message string, fields ...zap.Field) {
	l.Logger.Info(message, fields...)
}

// Warn -.
func (l *Logger) Warn(message string, fields ...zap.Field) {
	l.Logger.Warn(message, fields...)
}

// Error -.
func (l *Logger) Error(message string, fields ...zap.Field) {
	l.Logger.Error(message, fields...)
}

// Fatal -.
func (l *Logger) Fatal(message string, fields ...zap.Field) {
	l.Logger.Fatal(message, fields...)
}
