package properties

// Friend vault
type Friendship struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account string `json:"account" gorm:"not null"`
	Hash    string `json:"hash" gorm:"not null"`
	Payload string `json:"friend" gorm:"not null"` // Encrypted (with account's public key) friend key + data
}

// Vault for all kinds of things (e.g. conversation tokens, etc.)
type VaultEntry struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account string `json:"account" gorm:"not null"`
	Hash    string `json:"hash" gorm:"not null"`
	Payload string `json:"payload" gorm:"not null"` // Encrypted (with account's public key) data
}
