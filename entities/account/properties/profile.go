package properties

type Profile struct {
	ID string `json:"id" gorm:"primaryKey"` // Account ID

	Picture string `json:"picture"` // Picture
	Data    string `json:"data"`    // Encrypted data (if we need it in the future)
}
