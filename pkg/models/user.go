package models

type User struct {
	ID       string `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
}
