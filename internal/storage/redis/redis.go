package redis

import (
	"context"
	"fmt"
	"moscow-events-telegramauth/internal/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
}

func NewRedisClient(cfg *config.Config) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DbName,
	})

	return &Redis{rdb: rdb}
}

func (s *Redis) Ping(ctx context.Context) error {
	const op = "storage.redis.Ping"
	if err := s.rdb.Ping(ctx); err.Err() != nil {
		return fmt.Errorf("%s:%w", op, err.Err())
	}
	return nil
}
