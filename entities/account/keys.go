package account

type PublicKey struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Key       string `json:"key"`
	Signature string `json:"signature"`
}

type PrivateKey struct {
	ID  uint   `json:"id" gorm:"primaryKey"`
	Key string `json:"key"`
}
