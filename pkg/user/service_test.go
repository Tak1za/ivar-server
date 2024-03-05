package user

import (
	"errors"
	"ivar/pkg/database"
	"testing"
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
