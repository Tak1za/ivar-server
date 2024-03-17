package chat

import (
	"errors"
	"ivar/pkg/database"
	"ivar/pkg/models"
	"testing"
)

func TestService_AddMessage_Success(t *testing.T) {
	m := new(database.MockStore)
	m.On("StoreMessage", models.Message{
		Sender:    "senderId",
		Recipient: "recipientId",
		Content:   "test message 1",
	}).Return(nil)

	s := Service{m}

	err := s.AddMessage(models.Message{
		Sender:    "senderId",
		Recipient: "recipientId",
		Content:   "test message 1",
	})

	m.AssertExpectations(t)

	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestService_AddMessage_Failure(t *testing.T) {
	m := new(database.MockStore)
	m.On("StoreMessage", models.Message{
		Sender:    "senderId",
		Recipient: "recipientId",
		Content:   "test message 1",
	}).Return(errors.New("failed"))

	s := Service{m}

	err := s.AddMessage(models.Message{
		Sender:    "senderId",
		Recipient: "recipientId",
		Content:   "test message 1",
	})

	m.AssertExpectations(t)

	if err.Error() != "failed" {
		t.Errorf("error should be 'failed', got: %v", err)
	}
}
