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

func TestService_GetUser_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetUser", "username1").Return("userid1", "username1", nil)

	s := Service{m}

	user, err := s.Get("username1")

	m.AssertExpectations(t)

	if user.ID != "userid1" || user.Username != "username1" {
		t.Errorf("id should be %v, got: %v; username should be %v, got: %v", "userid1", user.ID, "username1", user.Username)
	}

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_GetUser_Failure(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetUser", "username1").Return("", "", errors.New("failed"))

	s := Service{m}

	user, err := s.Get("username1")

	m.AssertExpectations(t)

	if user != nil {
		t.Errorf("user should be nil, got: %v", user)
	}

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_AddFriendRequest_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetUser", "username1").Return("user1", "", nil)
	m.On("GetUser", "username2").Return("user2", "", nil)
	m.On("AddFriendRequest", "user1", "user2").Return(nil)

	s := Service{m}

	err := s.AddFriend(&models.AddFriendRequest{
		UserA: "username1",
		UserB: "username2",
	})

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_AddFriendRequest_Failure_FetchingUserA(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetUser", "username1").Return("", "", errors.New("failed"))

	s := Service{m}

	err := s.AddFriend(&models.AddFriendRequest{
		UserA: "username1",
		UserB: "username2",
	})

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_AddFriendRequest_Failure_FetchingUserB(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetUser", "username1").Return("user1", "", nil)
	m.On("GetUser", "username2").Return("", "", errors.New("failed"))

	s := Service{m}

	err := s.AddFriend(&models.AddFriendRequest{
		UserA: "username1",
		UserB: "username2",
	})

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_AddFriendRequest_Failure_SendingFriendRequest(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetUser", "username1").Return("user1", "", nil)
	m.On("GetUser", "username2").Return("user2", "", nil)
	m.On("AddFriendRequest", "user1", "user2").Return(errors.New("failed"))

	s := Service{m}

	err := s.AddFriend(&models.AddFriendRequest{
		UserA: "username1",
		UserB: "username2",
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
	m.On("GetUser", "username1").Return("user1", "", nil)
	m.On("GetFriendRequests", "user1").Return([]models.FriendRequest{{
		ID:    1,
		UserA: "user1",
		UserB: "user2",
	},
	}, nil)

	s := Service{m}

	_, err := s.GetFriendRequests("username1")

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_GetFriendrequests_Failure_FetchingUser(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetUser", "username1").Return("", "", errors.New("failed"))

	s := Service{m}

	_, err := s.GetFriendRequests("username1")

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_GetFriendrequests_Failure_GettingFriendRequests(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetUser", "username1").Return("user1", "", nil)
	m.On("GetFriendRequests", mock.AnythingOfType("string")).Return([]models.FriendRequest{}, errors.New("failed"))

	s := Service{m}

	_, err := s.GetFriendRequests("username1")

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}
