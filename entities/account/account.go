package account

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model

	Username string `json:"username"`
	Password string `json:"password"`
}