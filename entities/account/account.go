package account

import (
	"node-backend/util/auth"
	"time"
)

type Account struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Username  string    `json:"username"`
	Tag       string    `json:"tag" gorm:"unique"`
	Password  string    `json:"password"`
	Email     string    `json:"email" gorm:"unique"`
	RankID    uint      `json:"rank"`
	CreatedAt time.Time `json:"created_at"`

	Rank Rank `gorm:"foreignKey:RankID"`
	/*
		Sessions       []Session                   `gorm:"foreignKey:Account"`
		Subscription   Subscription                `gorm:"foreignKey:Account"`
		Authentication []Authentication            `gorm:"foreignKey:Account"`
		Friends        []properties.Friend         `gorm:"foreignKey:Account"`
		Settings       []properties.AccountSetting `gorm:"foreignKey:Account"`
	*/
}

func (a *Account) CheckPassword(password string) bool {

	// Check if password is correct
	return a.Password == auth.HashPassword(password)
}
