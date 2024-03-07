package models

type AddFriendRequest struct {
	UserA string `json:"userA"`
	UserB string `json:"userB"`
}

type UpdateFriendRequest struct {
	ID     int `json:"id"`
	Status int `json:"status"`
}

type FriendRequest struct {
	ID     int    `json:"id"`
	UserA  string `json:"userA"`
	UserB  string `json:"userB"`
	Status int    `json:"status"`
}
