package models

type Message struct {
	ID        int64  `json:"id"`
	Timestamp string `json:"timestamp"`
	Sender    string `json:"sender" binding:"required"`
	Recipient string `json:"recipient" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
