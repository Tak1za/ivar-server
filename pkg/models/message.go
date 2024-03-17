package models

type Message struct {
	ID        int64  `json:"id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Sender    string `json:"sender" binding:"required"`
	Recipient string `json:"recipient" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
