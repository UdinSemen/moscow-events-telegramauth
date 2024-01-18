package storage

import "context"

type RedisStorage interface {
	ConfirmSession(ctx context.Context, timeCode, chatID string) error
}
