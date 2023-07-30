package properties

type Friendship struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account string `json:"account" gorm:"not null"`
	Hash    string `json:"hash" gorm:"not null"`
	Payload string `json:"friend" gorm:"not null"` // Encrypted (with account's public key) friend key + data
}
