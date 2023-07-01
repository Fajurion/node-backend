package account

import (
	"regexp"
	"strings"
	"time"
)

type Account struct {
	ID string `json:"id" gorm:"primaryKey"` // 8 character-long string

	Email     string    `json:"email" gorm:"uniqueIndex"`
	Username  string    `json:"username"`
	Tag       string    `json:"tag"`
	RankID    uint      `json:"rank"`
	CreatedAt time.Time `json:"created_at"`

	Rank           Rank             `json:"-" gorm:"foreignKey:RankID"`
	Authentication []Authentication `json:"-" gorm:"foreignKey:Account"`
	Sessions       []Session        `json:"-" gorm:"foreignKey:Account"`
}

const EmailRegex = "^[a-zA-Z0-9]+(?:\\.[a-zA-Z0-9]+)*@[a-zA-Z0-9]+(?:\\.[a-zA-Z0-9]+)*$"

func NormalizeEmail(email string) string {

	// Convert email to lowercase
	email = strings.ToLower(email)

	// Remove leading and trailing whitespaces
	email = strings.TrimSpace(email)

	// Remove dots (.) from the username part of the email
	parts := strings.Split(email, "@")
	username := parts[0]
	username = strings.ReplaceAll(username, ".", "")

	// Reconstruct the normalized email address
	normalizedEmail := username + "@" + parts[1]

	return normalizedEmail
}

func CheckEmail(email string) (bool, string) {

	// Check if email is valid
	match, err := regexp.Match(EmailRegex, []byte(email))
	if !match || err != nil {
		return false, ""
	}

	email = NormalizeEmail(email)
	if strings.Contains(email, " ") {
		return false, ""
	}

	return true, email
}
