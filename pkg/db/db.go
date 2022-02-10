package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func NewDB(url string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
