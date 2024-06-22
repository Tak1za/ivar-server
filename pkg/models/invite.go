package models

import "time"

type Invite struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Code      string    `json:"code"`
	ServerId  int64     `json:"serverId"`
}
