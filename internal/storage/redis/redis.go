package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/UdinSemen/moscow-events-telegramauth/internal/config"
	models "github.com/UdinSemen/moscow-events-telegramauth/internal/models/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	timeCodeTTL   = "90s" //seconds
	timeCodeTable = "time_codes."
)

var (
	ErrEmptyJSONValue          = errors.New("empty value")
	ErrSessionAlreadyConfirmed = errors.New("session already confirmed")
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

func (s *Redis) ConfirmSession(ctx context.Context, timeCode, userID string) error {
	const op = "storage.redis.ConfirmSession"

	res, err := s.rdb.Get(ctx, timeCodeTable+timeCode).Result()
	if err != nil {
		zap.S().Debug(err)
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%s:%w", op, ErrEmptyJSONValue)
		}
		return fmt.Errorf("%s:%w", op, err)
	}
	zap.S().Debug(res)

	var session models.RegSession
	if err := json.Unmarshal([]byte(res), &session); err != nil {
		zap.S().Debug(res)
		return fmt.Errorf("%s:%w", op, err)
	}
	if session.IsConfirmed {
		return fmt.Errorf("%s:%w", op, ErrSessionAlreadyConfirmed)
	}
	session.IsConfirmed = true
	session.UserID = userID

	confirmedSession, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	dur, err := time.ParseDuration(timeCodeTTL)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	if err := s.rdb.Set(ctx, timeCodeTable+timeCode, confirmedSession, dur).Err(); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}
