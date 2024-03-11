package models

type AddFriendRequest struct {
	UsernameA string `json:"usernameA"`
	UsernameB string `json:"usernameB"`
}

type UpdateFriendRequest struct {
	ID     int `json:"id"`
	Status int `json:"status"`
}

type FriendRequest struct {
	ID     int  `json:"id"`
	UserA  User `json:"userA"`
	UserB  User `json:"userB"`
	Status int  `json:"status"`
}

type RemoveFriendRequest struct {
	CurrentUserId  string `json:"currentUserId"`
	ToRemoveUserId string `json:"toRemoveUserId"`
}
