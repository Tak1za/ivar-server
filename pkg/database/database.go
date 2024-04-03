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
	AddFriendRequest(userIdA, userIdB string) error
	UpdateFriendRequest(id, status int) error
	GetFriendRequests(currentUserId string) ([]models.FriendRequest, error)
	GetFriends(username string) ([]models.User, error)
	RemoveFriend(currentUserId, toRemoveUserId string) error
	GetChatInfo(users []string) (models.ChatInfo, error)
	StoreMessage(message models.Message) error
	RetrieveMessages(users []string) ([]models.Message, error)
	AllChats(userId string) ([]models.User, error)
	CreateServer(name, userId string) error
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{
		db: db,
	}
}

func (s *store) CreateServer(name, userId string) error {
	query := "insert into servers (name, owner) values (@name, @owner) on conflict do nothing"
	args := pgx.NamedArgs{
		"name":  name,
		"owner": userId,
	}
	_, err := s.db.Exec(context.Background(), query, args)
	if err != nil {
		log.Println("unable to insert row: " + err.Error())
		return err
	}

	return nil
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

func (s *store) AddFriendRequest(userIdA, userIdB string) error {
	query := "insert into friends (user_a, user_b) values (@userA, @userB)"
	args := pgx.NamedArgs{
		"userA": userIdA,
		"userB": userIdB,
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

func (s *store) GetFriendRequests(currentUserId string) ([]models.FriendRequest, error) {
	rows, _ := s.db.Query(context.Background(), `select f.id, jsonb_build_object('id', fromUser.id, 'username', fromUser.username) as userA, jsonb_build_object('id', toUser.id, 'username', toUser.username) as userB, f.status
	from friends f
	inner join users fromUser on
	f.user_a = fromUser.id
	inner join users toUser on
	f.user_b = toUser.id
	where (fromUser.id = $1 or toUser.id = $1) and status = 2`, currentUserId)
	friendRequests, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.FriendRequest])
	if err != nil {
		log.Println("unable to fetch rows: " + err.Error())
		return []models.FriendRequest{}, err
	}

	return friendRequests, nil
}

func (s *store) GetFriends(userId string) ([]models.User, error) {
	rows, _ := s.db.Query(context.Background(), `select json_agg(distinct jsonb_build_object('id', u.id, 'username', u.username)) from friends f
	inner join users u
	on u.id = f.user_a or u.id = f.user_b
	where ((f.user_a = $1 or f.user_b = $1) and f.status = 1)`, userId)
	friends, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) ([]models.User, error) {
		res := make([]models.User, 0)
		err := row.Scan(&res)
		return res, err
	})
	if err != nil {
		log.Println("unable to fetch rows: ", err.Error())
		return nil, err
	}

	return friends, nil
}

func (s *store) RemoveFriend(currentUserId, toRemoveUserId string) error {
	query := `delete from
	friends as f
  	where
	(
	  f.user_a = @userA
	  and 
	  f.user_b = @userB
	)
	or
	(
	  f.user_a = @userB
	  and 
	  f.user_b = @userA
	)`
	args := pgx.NamedArgs{
		"userA": currentUserId,
		"userB": toRemoveUserId,
	}

	if _, err := s.db.Exec(context.Background(), query, args); err != nil {
		log.Println("unable to delete row: " + err.Error())
		return err
	}

	return nil
}

func (s *store) GetChatInfo(users []string) (models.ChatInfo, error) {
	rows, _ := s.db.Query(context.Background(), "select id, username from users where id = any($1)", users)

	userDetails, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		log.Println("unable to fetch rows: " + err.Error())
		return models.ChatInfo{}, err
	}

	return models.ChatInfo{Users: userDetails}, nil
}

func (s *store) StoreMessage(message models.Message) error {
	query := "insert into messages (sender_id, recipient_id, content) values (@senderId, @recipientId, @content)"
	args := pgx.NamedArgs{
		"senderId":    message.Sender,
		"recipientId": message.Recipient,
		"content":     message.Content,
	}

	if _, err := s.db.Exec(context.Background(), query, args); err != nil {
		log.Println("unable to insert row: " + err.Error())
		return err
	}

	return nil
}

func (s *store) RetrieveMessages(users []string) ([]models.Message, error) {
	rows, _ := s.db.Query(context.Background(), `select id, timestamp, content, sender_id as sender, recipient_id as recipient from messages m where (sender_id = $1 or recipient_id = $1) and (sender_id = $2 or recipient_id = $2) order by timestamp desc`, users[0], users[1])
	messages, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Message])
	if err != nil {
		log.Println("unable to fetch rows: " + err.Error())
		return []models.Message{}, err
	}

	return messages, nil
}

func (s *store) AllChats(userId string) ([]models.User, error) {
	rows, _ := s.db.Query(context.Background(), `select
	json_agg(
	  	distinct jsonb_build_object('id', u.id, 'username', u.username)
	)
  	from
		messages m
		inner join 
		users u
		on 
		u.id = m.sender_id 
		or
		u.id = m.recipient_id
  	where (m.sender_id = $1 or m.recipient_id = $1) and u.id <> $1`, userId)
	chats, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) ([]models.User, error) {
		res := make([]models.User, 0)
		err := row.Scan(&res)
		return res, err
	})
	if err != nil {
		log.Println("unable to fetch rows: ", err.Error())
		return nil, err
	}

	return chats, nil
}
