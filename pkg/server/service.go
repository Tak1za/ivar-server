package server

import (
	"crypto/rand"
	"ivar/pkg/database"
	"ivar/pkg/models"
	"log"
	"math/big"
)

type Service struct {
	Store database.Store
}

const (
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	inviteLength = 6
)

func (s *Service) CreateServer(name, userId string) error {
	if err := s.Store.CreateServer(name, userId); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetServers() ([]models.Server, error) {
	return s.Store.GetServers()
}

func (s *Service) CreateInvite(serverId int) (string, error) {
	existingCode, err := s.Store.GetInvite(serverId)
	if err != nil {
		log.Println("unable to get invite: " + err.Error())
		return "", err
	}

	if existingCode != "" {
		return existingCode, nil
	}

	// generate a code
	code := make([]byte, inviteLength)
	for i := range code {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[randomIndex.Int64()]
	}

	if err := s.Store.StoreInvite(string(code), serverId); err != nil {
		if err.Error() == "invalid server" {
			log.Println("invalid server")
			return "", err
		}
		log.Println("error creating invite: " + err.Error())
		return "", err
	}

	return string(code), nil
}
