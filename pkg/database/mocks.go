package database

import "github.com/stretchr/testify/mock"

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateUser(id, username string) error {
	returnVals := m.Called(id, username)

	return returnVals.Error(0)
}
