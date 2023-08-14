package properties

type StoredAction struct {
	ID string `json:"-" gorm:"primaryKey"`

	Account string `json:"-" gorm:"not null"`
	Payload string `json:"action" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
}

// Authenticated stored actions
type AStoredAction struct {
	ID string `json:"-" gorm:"primaryKey"`

	Account string `json:"-" gorm:"not null"`
	Payload string `json:"action" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
}
