package user

import "ivar/pkg/database"

type Service struct {
	Store database.Store
}

func (s *Service) Create(id, username string) error {
	if err := s.Store.CreateUser(id, username); err != nil {
		return err
	}
	return nil
}
