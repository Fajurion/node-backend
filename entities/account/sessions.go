package account

import "time"

type Session struct {
	token string `json:"token" gorm:"primaryKey"`

	Account     uint      `json:"account"`
	AccountName string    `json:"account_name"`
	Device      string    `json:"device"`
	Connected   bool      `json:"connected"`
	End         time.Time `json:"end"`
	App         uint      `json:"app"`
	Node        uint      `json:"node"`
}
