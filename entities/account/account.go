package account

import (
	"node-backend/util/auth"
	"strings"
	"time"
)

type Account struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Username  string    `json:"username"`
	Tag       string    `json:"tag"`
	Password  string    `json:"password"`
	Email     string    `json:"email" gorm:"unique"`
	RankID    uint      `json:"rank"`
	CreatedAt time.Time `json:"created_at"`

	Rank           Rank             `json:"-" gorm:"foreignKey:RankID"`
	Authentication []Authentication `json:"-" gorm:"foreignKey:Account"`
	Sessions       []Session        `json:"-" gorm:"foreignKey:Account"`
	Subscription   Subscription     `json:"-" gorm:"foreignKey:Account"`
	Key            PublicKey        `json:"key" gorm:"foreignKey:ID"`
}

func (a *Account) CheckPassword(password string) bool {

	// Check if password is correct
	return strings.Compare(a.Password, auth.HashPassword(password)) == 0
}
