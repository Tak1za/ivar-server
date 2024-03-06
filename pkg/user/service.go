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

func (s *Service) AddFriend(request *models.FriendRequest) error {
	if err := s.Store.AddFriendRequest(request.UserA, request.UserB); err != nil {
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
