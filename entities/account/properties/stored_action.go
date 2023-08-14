package properties

type StoredAction struct {
	ID string `json:"-" gorm:"primaryKey"`

	Account string `json:"-" gorm:"not null"`
	Payload string `json:"action" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
}

// TODO: Implement some sort of key that one can send to friends for them to be able to send those actions
type AuthenticatedStoredAction struct {
	ID string `json:"-" gorm:"primaryKey"`

	Account string `json:"-" gorm:"not null"`
	Payload string `json:"action" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
}
