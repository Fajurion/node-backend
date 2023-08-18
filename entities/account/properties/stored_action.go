package properties

type StoredAction struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account string `json:"-" gorm:"not null"`
	Payload string `json:"payload" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
}

// Authenticated stored actions
type AStoredAction struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account string `json:"-" gorm:"not null"`
	Payload string `json:"payload" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
}
