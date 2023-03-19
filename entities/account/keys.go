package account

type PublicKey struct {
	ID  uint   `json:"id" gorm:"primaryKey"`
	Key string `json:"key"`
}

type PrivateKey struct {
	ID  uint   `json:"id" gorm:"primaryKey"`
	Key string `json:"key"`
}
