package main

import (
	"context"

	"github.com/UdinSemen/moscow-events-telegramauth/internal/bot/telegram"
	"github.com/UdinSemen/moscow-events-telegramauth/internal/config"
	pg_storage "github.com/UdinSemen/moscow-events-telegramauth/internal/storage/pg-storage"
	"github.com/UdinSemen/moscow-events-telegramauth/internal/storage/redis"
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

	redisStorage := redis.NewRedisClient(cfg)
	if err := redisStorage.Ping(context.Background()); err != nil {
		zap.S().Fatal(err.Error())
	}
	postgresStorage, err := pg_storage.InitPgStorage(cfg)
	if err := postgresStorage.Ping(); err != nil {
		zap.S().Fatalf(err.Error())
	}

	bot, err := telegram.InitBot(cfg, redisStorage, postgresStorage, false)
	if err != nil {
		zap.S().Fatal(err.Error())
	}

	bot.Start()
}
