package server

import (
	"ivar/pkg/database"
	"ivar/pkg/models"
)

type Service struct {
	Store database.Store
}

func (s *Service) CreateServer(name, userId string) error {
	if err := s.Store.CreateServer(name, userId); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetServers() ([]models.Server, error) {
	return s.Store.GetServers()
}
