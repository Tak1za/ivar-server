package models

type ChatInfo struct {
	Users    []User    `json:"users"`
	Messages []Message `json:"messages"`
}

type ChatInfoRequest struct {
	Users []string `json:"users" binding:"required"`
}
