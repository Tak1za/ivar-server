package models

import "time"

type Server struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	OwnerID   string    `json:"ownerId" binding:"required"`
}
