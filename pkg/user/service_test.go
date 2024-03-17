package user

import (
	"errors"
	"ivar/pkg/database"
	"ivar/pkg/models"
	"testing"
	"time"

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
	m.On("GetUser", "username1").Return("user1", "", nil)
	m.On("GetUser", "username2").Return("user2", "", nil)
	m.On("AddFriendRequest", "user1", "user2").Return(nil)

	s := Service{m}

	err := s.AddFriend(&models.AddFriendRequest{
		UsernameA: "username1",
		UsernameB: "username2",
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
		UsernameA: "username1",
		UsernameB: "username2",
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
		UsernameA: "username1",
		UsernameB: "username2",
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
		UsernameA: "username1",
		UsernameB: "username2",
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

func TestService_GetFriends_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetFriends", "userId1").Return([]models.User{
		{ID: "userId1", Username: "username1"},
		{ID: "userId2", Username: "username2"},
	}, nil)

	s := Service{m}

	_, err := s.GetFriends("userId1")

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_GetFriends_NoCurrentUser_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetFriends", "userId1").Return([]models.User{
		{ID: "userId2", Username: "username2"},
	}, nil)

	s := Service{m}

	_, err := s.GetFriends("userId1")

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_GetFriends_Failure(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetFriends", "userId1").Return([]models.User{}, errors.New("failed"))

	s := Service{m}

	_, err := s.GetFriends("userId1")

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_RemoveFriend_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("RemoveFriend", "userId1", "userId2").Return(nil)

	s := Service{m}

	err := s.RemoveFriend("userId1", "userId2")

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_RemoveFriend_Failure(t *testing.T) {
	m := new(database.MockStore)
	m.On("RemoveFriend", "userId1", "userId2").Return(errors.New("failed"))

	s := Service{m}

	err := s.RemoveFriend("userId1", "userId2")

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_GetChatInfo_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetChatInfo", []string{"user1", "user2"}).Return(models.ChatInfo{Users: []models.User{
		{ID: "user1", Username: "username1"},
		{ID: "user2", Username: "username2"},
	}}, nil)
	m.On("RetrieveMessages", []string{"user1", "user2"}).Return([]models.Message{{ID: 1, Timestamp: time.Now(), Sender: "user1", Recipient: "user2"}}, nil)

	s := Service{m}

	_, err := s.GetChatInfo([]string{"user1", "user2"})

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_GetChatInfo_Failure(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetChatInfo", []string{"user1", "user2"}).Return(models.ChatInfo{}, errors.New("failed"))
	s := Service{m}

	_, err := s.GetChatInfo([]string{"user1", "user2"})

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}

func TestService_GetChatInfo_FailsToFetchMessages_Failure(t *testing.T) {
	m := new(database.MockStore)
	m.On("GetChatInfo", []string{"user1", "user2"}).Return(models.ChatInfo{Users: []models.User{
		{ID: "user1", Username: "username1"},
		{ID: "user2", Username: "username2"},
	}}, nil)
	m.On("RetrieveMessages", []string{"user1", "user2"}).Return([]models.Message{}, errors.New("failed"))
	s := Service{m}

	_, err := s.GetChatInfo([]string{"user1", "user2"})

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}
