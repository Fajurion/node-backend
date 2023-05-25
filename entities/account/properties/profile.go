package properties

type Profile struct {
	ID string `json:"id" gorm:"primaryKey"` // Account ID

	Key  string `json:"key"`  // AES key encrypted with the account's public key
	Data string `json:"data"` // AES encrypted data
}
