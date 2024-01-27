package pg_storage

import (
	"errors"
	"fmt"

	"github.com/UdinSemen/moscow-events-telegramauth/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	ErrPgUniqueConstrCode = "23505"
)

var (
	ErrPgUniqueConstr = errors.New("err postgres unique constraint")
)

type PgStorage struct {
	db *sqlx.DB
}

func InitPgStorage(cfg *config.Config) (*PgStorage, error) {
	const op = "pg_storage.InitPgStorage"

	dbConf := cfg.Postgres

	connect, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DbName, dbConf.SslMode))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &PgStorage{
		db: connect,
	}, nil
}

func (s *PgStorage) Ping() error {
	return s.db.Ping()
}

func (s *PgStorage) AddUser(firstName, lastName, sex string, userID int64) error {
	const op = "pg_storage.AddUser"

	_, err := s.db.NamedExec("insert into users(tg_user_id, first_name, last_name, sex) "+
		"values (:tg_user_id, :fir_name, :last_name, :sex)", map[string]interface{}{
		"tg_user_id": userID,
		"fir_name":   firstName,
		"last_name":  lastName,
		"sex":        sex,
	})

	if err != nil {
		if err.(*pq.Error).Code == ErrPgUniqueConstrCode {
			return fmt.Errorf("%s: %w", op, ErrPgUniqueConstr)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
