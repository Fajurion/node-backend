package account

type PublicKey struct {
	ID  uint   `json:"id" gorm:"primaryKey"`
	Key string `json:"key"`
}

// TODO: Private key?
