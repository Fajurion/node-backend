package account

type PublicKey struct {
	ID  string `json:"id" gorm:"primaryKey"`
	Key string `json:"key"`
}

type PrivateKey struct {
	ID  string `json:"id" gorm:"primaryKey"`
	Key string `json:"key"`
}
