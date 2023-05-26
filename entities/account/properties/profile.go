package properties

// * This object is completely exposed to all friends, BE CAREFUL CHANGING IT
type Profile struct {
	ID string `json:"id" gorm:"primaryKey"` // Account ID

	Key  string `json:"key"`  // AES key encrypted with the account's public key
	Data string `json:"data"` // AES encrypted data
}
