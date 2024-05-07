package database

import (
	"ivar/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateUser(id, username string) error {
	returnVals := m.Called(id, username)

	return returnVals.Error(0)
}

func (m *MockStore) GetUser(username string) (string, string, error) {
	returnVals := m.Called(username)

	return returnVals.String(0), returnVals.String(1), returnVals.Error(2)
}

func (m *MockStore) AddFriendRequest(userA, userB string) error {
	returnVals := m.Called(userA, userB)

	return returnVals.Error(0)
}

func (m *MockStore) UpdateFriendRequest(id, status int) error {
	returnVals := m.Called(id, status)

	return returnVals.Error(0)
}

func (m *MockStore) GetFriendRequests(userA string) ([]models.FriendRequest, error) {
	returnVals := m.Called(userA)

	return returnVals.Get(0).([]models.FriendRequest), returnVals.Error(1)
}

func (m *MockStore) GetFriends(userId string) ([]models.User, error) {
	returnVals := m.Called(userId)

	return returnVals.Get(0).([]models.User), returnVals.Error(1)
}

func (m *MockStore) RemoveFriend(usernameA, usernameB string) error {
	returnVals := m.Called(usernameA, usernameB)

	return returnVals.Error(0)
}

func (m *MockStore) GetChatInfo(users []string) (models.ChatInfo, error) {
	returnVals := m.Called(users)

	return returnVals.Get(0).(models.ChatInfo), returnVals.Error(1)
}

func (m *MockStore) StoreMessage(message models.Message) error {
	returnVals := m.Called(message)

	return returnVals.Error(0)
}

func (m *MockStore) RetrieveMessages(users []string) ([]models.Message, error) {
	returnVals := m.Called(users)

	return returnVals.Get(0).([]models.Message), returnVals.Error(1)
}

func (m *MockStore) AllChats(userId string) ([]models.User, error) {
	returnVals := m.Called(userId)

	return returnVals.Get(0).([]models.User), returnVals.Error(1)
}

func (m *MockStore) CreateServer(name, userId string) error {
	returnVals := m.Called(name, userId)

	return returnVals.Error(0)
}

func (m *MockStore) GetServers() ([]models.Server, error) {
	returnVals := m.Called()

	return returnVals.Get(0).([]models.Server), returnVals.Error(1)
}
