package account

import (
	"time"
)

type Session struct {
	ID    uint   `json:"id" gorm:"primaryKey"` // ID is the primary key of the table
	Token string `json:"token" gorm:"unique"`

	Account         uint      `json:"account"`
	PermissionLevel uint      `json:"level"`
	Device          string    `json:"device"`
	Connected       bool      `json:"connected"`
	End             time.Time `json:"end"`
	App             uint      `json:"app"`
	Node            uint      `json:"node"`
}

func (s *Session) IsExpired() bool {
	return time.Since(s.End) > 0
}
