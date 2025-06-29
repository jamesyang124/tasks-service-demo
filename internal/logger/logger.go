package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Package logger provides a singleton logger instance using Uber's zap.

var (
	// once ensures the logger is only initialized once.
	once sync.Once
	// instance holds the singleton zap.SugaredLogger.
	instance *zap.SugaredLogger
)

// Get returns a singleton SugaredLogger instance for application-wide logging.
// The logger is configured with production settings and ISO8601 time encoding.
func Get() *zap.SugaredLogger {
	once.Do(func() {
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "ts"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		log, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		instance = log.Sugar()
	})
	return instance
}
