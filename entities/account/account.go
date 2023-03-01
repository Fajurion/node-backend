package account

import (
	"node-backend/entities/account/properties"
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

	Rank           Rank                        `gorm:"foreignKey:RankID"`
	Authentication []Authentication            `gorm:"foreignKey:Account"`
	Sessions       []Session                   `gorm:"foreignKey:Account"`
	Subscription   Subscription                `gorm:"foreignKey:Account"`
	Friends        []properties.Friend         `gorm:"foreignKey:Account"`
	Settings       []properties.AccountSetting `gorm:"foreignKey:Account"`
}

func (a *Account) CheckPassword(password string) bool {

	// Check if password is correct
	return strings.Compare(a.Password, auth.HashPassword(password)) == 0
}
