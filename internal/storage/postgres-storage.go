package storage

type PgStorage interface {
	Ping() error
	AddUser(firstName, lastName, sex string, userID int64) error
}
