package models

type ChatInfo struct {
	Users []User `json:"users"`
}

type ChatInfoRequest struct {
	Users []string `json:"users" binding:"required"`
}
