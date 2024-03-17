package models

import "time"

type Message struct {
	ID        int64     `json:"id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Sender    string    `json:"sender" binding:"required"`
	Recipient string    `json:"recipient" binding:"required"`
	Content   string    `json:"content" binding:"required"`
}
