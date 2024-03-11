package user

import (
	"errors"
	"ivar/pkg/database"
	"ivar/pkg/models"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestService_CreateUser_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("CreateUser", "testid1", "testusername1").Return(nil)

	s := Service{m}

	err := s.Create("testid1", "testusername1")

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_CreateUser_Failure(t *testing.T) {
	m := new(database.MockStore)
	m.On("CreateUser", "testid1", "testusername1").Return(errors.New("failed"))

	s := Service{m}

	err := s.Create("testid1", "testusername1")

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_AddFriendRequest_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("AddFriendRequest", "user1", "user2").Return(nil)

	s := Service{m}

	err := s.AddFriend(&models.AddFriendRequest{
		UsernameA: "user1",
		UsernameB: "user2",
	})

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_AddFriendRequest_Failure(t *testing.T) {
	m := new(database.MockStore)
	m.On("AddFriendRequest", "user1", "user2").Return(errors.New("failed"))

	s := Service{m}

	err := s.AddFriend(&models.AddFriendRequest{
		UsernameA: "user1",
		UsernameB: "user2",
	})

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_UpdateFriendRequest_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("UpdateFriendRequest", 1, 1).Return(nil)

	s := Service{m}

	err := s.UpdateFriend(&models.UpdateFriendRequest{
		ID:     1,
		Status: 1,
	})

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_UpdateFriendRequest_Failure(t *testing.T) {
	m := new(database.MockStore)
	m.On("UpdateFriendRequest", 1, 1).Return(errors.New("failed"))

	s := Service{m}

	err := s.UpdateFriend(&models.UpdateFriendRequest{
		ID:     1,
		Status: 1,
	})

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_GetFriendrequests_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetFriendRequests", "1").Return([]models.FriendRequest{{
		ID: 1,
		UserA: models.User{
			ID:       "1",
			Username: "user1",
		},
		UserB: models.User{
			ID:       "2",
			Username: "user2",
		},
	},
	}, nil)

	s := Service{m}

	_, err := s.GetFriendRequests("1")

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_GetFriendrequests_Failure_GettingFriendRequests(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetFriendRequests", mock.AnythingOfType("string")).Return([]models.FriendRequest{}, errors.New("failed"))

	s := Service{m}

	_, err := s.GetFriendRequests("1")

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}
