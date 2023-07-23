package properties

type Friendship struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Account string `json:"account" gorm:"primaryKey"`
	Payload string `json:"friend"` // Encrypted (with account's public key) friend key + data
}
