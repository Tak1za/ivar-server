package models

type AddFriendRequest struct {
	UsernameA string `json:"usernameA" binding:"required"`
	UsernameB string `json:"usernameB" binding:"required"`
}

type UpdateFriendRequest struct {
	ID     int `json:"id" binding:"required"`
	Status int `json:"status" binding:"required"`
}

type FriendRequest struct {
	ID     int  `json:"id"`
	UserA  User `json:"userA"`
	UserB  User `json:"userB"`
	Status int  `json:"status"`
}

type RemoveFriendRequest struct {
	CurrentUserId  string `json:"currentUserId" binding:"required"`
	ToRemoveUserId string `json:"toRemoveUserId" binding:"required"`
}
