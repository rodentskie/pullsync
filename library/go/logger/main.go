package logger

import (
	"slack-pr-lambda/env"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LoggerConfig() zap.Config {
	GOENV := env.GetEnv("ENV", "dev")

	config := zap.NewProductionConfig()

	if GOENV == "test" {
		// disable logs during test environment
		config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	}

	return config
}
