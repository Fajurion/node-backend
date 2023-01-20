package account

import "time"

type Session struct {
	Token string `json:"token" gorm:"primaryKey"`

	Account         uint      `json:"account"`
	PermissionLevel uint      `json:"level"`
	Device          string    `json:"device"`
	Connected       bool      `json:"connected"`
	End             time.Time `json:"end"`
	App             uint      `json:"app"`
	Node            uint      `json:"node"`
}
