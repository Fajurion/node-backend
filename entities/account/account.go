package account

import "time"

type Account struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Username  string    `json:"username"`
	Tag       string    `json:"tag" gorm:"unique"`
	Password  string    `json:"password"`
	Email     string    `json:"email" gorm:"unique"`
	Rank      uint      `json:"rank"`
	CreatedAt time.Time `json:"created_at"`
}

// TODO: Hash password
func (a *Account) CheckPassword(password string) bool {

	// Check if password is correct
	return a.Password == password
}
