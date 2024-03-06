package models

type FriendRequest struct {
	UserA string `json:"userA"`
	UserB string `json:"userB"`
}

type UpdateFriendRequest struct {
	ID     int `json:"id"`
	Status int `json:"status"`
}
