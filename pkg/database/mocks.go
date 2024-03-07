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
