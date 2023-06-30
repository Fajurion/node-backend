package account

import (
	"node-backend/util/auth"
	"strings"
)

type Authentication struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Account string `json:"account"`
	Type    uint   `json:"type"`
	Secret  string `json:"secret"`
}

const TypePassword = 0
const TypeTOTP = 1
const TypeRecoveryCode = 2
const TypePasskey = 3 // Implemented in the future

// Order to autenticate (0 = first, 1 = second, etc.)
var Order = map[uint]uint{
	TypePassword:     0,
	TypePasskey:      5, // Disabled (needs to still be implemented), will eventually be first too
	TypeTOTP:         1,
	TypeRecoveryCode: 1,
}

// Starting step when authenticating
const StartStep = 0

func (a *Authentication) checkPassword(password string) bool {
	return strings.Compare(a.Secret, auth.HashPassword(password)) == 0
}

func (a *Authentication) Verify(authType uint, secret string) bool {

	if a.Type != authType {
		return false
	}

	switch authType {
	case TypePassword:
		return a.checkPassword(secret)
	case TypeTOTP:
		return false // TODO: Implement
	case TypeRecoveryCode:
		return strings.Compare(a.Secret, secret) == 0 // TODO: Implement
	}

	return false
}
