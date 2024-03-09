package user

import (
	"ivar/pkg/database"
	"ivar/pkg/models"
)

type Service struct {
	Store database.Store
}

func (s *Service) Create(id, username string) error {
	if err := s.Store.CreateUser(id, username); err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(name string) (*models.User, error) {
	id, username, err := s.Store.GetUser(name)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       id,
		Username: username,
	}, nil
}

func (s *Service) AddFriend(request *models.AddFriendRequest) error {
	userIdA, _, err := s.Store.GetUser(request.UserA)
	if err != nil {
		return err
	}

	userIdB, _, err := s.Store.GetUser(request.UserB)
	if err != nil {
		return err
	}

	if err := s.Store.AddFriendRequest(userIdA, userIdB); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateFriend(request *models.UpdateFriendRequest) error {
	if err := s.Store.UpdateFriendRequest(request.ID, request.Status); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetFriendRequests(userA string) ([]models.FriendRequest, error) {
	friendRequests, err := s.Store.GetFriendRequests(userA)
	if err != nil {
		return nil, err
	}

	return friendRequests, nil
}

func (s *Service) GetFriends(username string) ([]string, error) {
	friends, err := s.Store.GetFriends(username)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func (s *Service) RemoveFriend(usernameA, usernameB string) error {
	if err := s.Store.RemoveFriend(usernameA, usernameB); err != nil {
		return err
	}

	return nil
}
