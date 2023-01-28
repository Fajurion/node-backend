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

func (s *Session) IsDesktop() bool {
	return s.Device == "desktop"
}

func (s *Session) IsWeb() bool {
	return s.Device == "web"
}

func (s *Session) IsExpired() bool {
	return time.Since(s.End) > 0
}

func (s *Session) Upgrade() {
	s.Device = "desktop"
}
