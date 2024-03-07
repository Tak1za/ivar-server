package database

import (
	"context"
	"ivar/pkg/models"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	CreateUser(id, username string) error
	GetUser(name string) (string, string, error)
	AddFriendRequest(userA, userB string) error
	UpdateFriendRequest(id, status int) error
	GetFriendRequests(userA string) ([]models.FriendRequest, error)
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
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

func (s *store) GetUser(name string) (string, string, error) {
	var (
		id       string
		username string
	)
	if err := s.db.QueryRow(context.Background(), "select id, username from users where username=$1", name).Scan(&id, &username); err != nil {
		log.Println("unable to query row: " + err.Error())
		return "", "", err
	}

	return id, username, nil
}

func (s *store) AddFriendRequest(userA, userB string) error {
	query := "insert into friends (user_a, user_b) values (@userA, @userB)"
	args := pgx.NamedArgs{
		"userA": userA,
		"userB": userB,
	}

	if _, err := s.db.Exec(context.Background(), query, args); err != nil {
		log.Println("unable to insert row: " + err.Error())
		return err
	}

	return nil
}

func (s *store) UpdateFriendRequest(id, status int) error {
	query := "update friends set status = @status where id = @id"
	args := pgx.NamedArgs{
		"status": status,
		"id":     id,
	}

	if _, err := s.db.Exec(context.Background(), query, args); err != nil {
		log.Println("unable to update row: " + err.Error())
		return err
	}

	return nil
}

func (s *store) GetFriendRequests(userA string) ([]models.FriendRequest, error) {
	rows, _ := s.db.Query(context.Background(), `select f.id, fromUser.username as userA, toUser.username as userB, f.status
	from friends f
	inner join users fromUser on
	f.user_a = fromUser.id
	inner join users toUser on
	f.user_b = toUser.id
	where (fromUser.username = $1 or toUser.username = $2) and status = 2`, userA, userA)
	friendRequests, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.FriendRequest])
	if err != nil {
		log.Println("unable to fetch rows: " + err.Error())
		return []models.FriendRequest{}, err
	}

	return friendRequests, nil
}
