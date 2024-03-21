package user

import (
	"ivar/pkg/database"
	"ivar/pkg/models"
	"slices"
	"strings"
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

func (s *Service) AddFriend(request *models.AddFriendRequest) error {
	userIdA, _, err := s.Store.GetUser(strings.ToLower(request.UsernameA))
	if err != nil {
		return err
	}

	userIdB, _, err := s.Store.GetUser(strings.ToLower(request.UsernameB))
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

func (s *Service) GetFriendRequests(currentUserId string) ([]models.FriendRequest, error) {
	friendRequests, err := s.Store.GetFriendRequests(currentUserId)
	if err != nil {
		return nil, err
	}

	return friendRequests, nil
}

func (s *Service) GetFriends(userId string) ([]models.User, error) {
	friends, err := s.Store.GetFriends(userId)
	if err != nil {
		return nil, err
	}

	currentUserIndex := slices.IndexFunc(friends, func(friend models.User) bool {
		return friend.ID == userId
	})

	if currentUserIndex == -1 {
		return friends, nil
	}

	before := friends[:currentUserIndex]
	after := friends[currentUserIndex+1:]

	return append(before, after...), nil
}

func (s *Service) RemoveFriend(currentUserId, toRemoveUserId string) error {
	if err := s.Store.RemoveFriend(currentUserId, toRemoveUserId); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetChatInfo(users []string) (models.ChatInfo, error) {
	chatInfo, err := s.Store.GetChatInfo(users)
	if err != nil {
		return models.ChatInfo{}, err
	}

	chatMessages, err := s.Store.RetrieveMessages(users)
	if err != nil {
		return models.ChatInfo{}, err
	}
	chatInfo.Messages = chatMessages

	return chatInfo, nil
}
