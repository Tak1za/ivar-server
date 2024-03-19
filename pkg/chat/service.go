package chat

import (
	"ivar/pkg/database"
	"ivar/pkg/models"
)

type Service struct {
	Store database.Store
}

func (s *Service) AddMessage(message models.Message) error {
	if err := s.Store.StoreMessage(message); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAllChats(userId string) ([]models.User, error) {
	users, err := s.Store.AllChats(userId)
	if err != nil {
		return []models.User{}, err
	}

	return users, nil
}
