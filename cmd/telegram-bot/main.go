package main

import (
	"github.com/UdinSemen/moscow-events-telegramauth/internal/config"
	"github.com/UdinSemen/moscow-events-telegramauth/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	dev   = "dev"
	local = "local"
	prod  = "prod"
)

func main() {
	cfg := config.MustLoad()

	var zapLevel zapcore.Level
	switch cfg.Env {
	case dev:
		zapLevel = zapcore.DebugLevel
	case local:
		zapLevel = zapcore.DebugLevel
	case prod:
		zapLevel = zapcore.WarnLevel
	}
	logger, err := utils.CreateLogger(zapLevel)
	defer func() {
		if err = logger.Sync(); err != nil {
			logger.Sugar().Error(err)
		}
	}()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
