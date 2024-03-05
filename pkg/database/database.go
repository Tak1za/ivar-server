package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type Store interface {
	CreateUser(id, username string) error
}

type store struct {
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) Store {
	return &store{
		db: db,
	}
}

func (s *store) CreateUser(id, username string) error {
	query := "insert into users (id, username) values (@id, @username) on conflict do nothing"
	args := pgx.NamedArgs{
		"id":       id,
		"username": username,
	}
	_, err := s.db.Exec(context.Background(), query, args)
	if err != nil {
		log.Println("unable to insert row: " + err.Error())
		return err
	}

	return nil
}
