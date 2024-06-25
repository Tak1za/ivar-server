package server

import (
	"crypto/rand"
	"ivar/pkg/database"
	"ivar/pkg/models"
	"log"
	"math/big"
	"time"
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
	invite, err := s.Store.GetInvite(serverId)
	if err != nil {
		log.Println("unable to get invite: " + err.Error())
		return "", err
	}

	if invite.Code != "" && invite.CreatedAt.Add(time.Hour*24).Compare(time.Now()) == 1 {
		log.Println("invite code is still valid for server: ", invite.ServerId)
		return invite.Code, nil
	}

	log.Printf("invite code does not exist or invite code created at %s for server %d is now invalid since it is over a day old", invite.CreatedAt.String(), invite.ServerId)

	if invite.Code != "" {
		if err := s.Store.DeleteInvite(serverId); err != nil {
			log.Println("error deleting invite: " + err.Error())
			return "", err
		}
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

func (s *Service) ValidateInvite(serverId int, code string) (bool, error) {
	invite, err := s.Store.GetInvite(serverId)
	if err != nil {
		log.Println("error getting invite code")
		return false, err
	}

	if invite.Code == "" {
		log.Println("no invite code found")
		return false, nil
	}

	if invite.CreatedAt.Add(time.Hour*24).Compare(time.Now()) == -1 {
		log.Printf("invite code created at %s is now invalid since it is over a day old", invite.CreatedAt.String())
		return false, nil
	}

	if invite.Code != code {
		log.Println("invalid code")
		return false, nil
	}

	return true, nil
}
