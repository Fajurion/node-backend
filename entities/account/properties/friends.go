package properties

import "node-backend/entities/account"

type Friend struct {
	Account   uint   `json:"account" gorm:"primaryKey"`
	Friend    uint   `json:"friend"`
	Updated   int64  `gorm:"autoUpdateTime:milli"`
	Request   bool   `json:"request"`
	Signature string `json:"signature"`

	AccountData account.Account   `json:"-" gorm:"foreignKey:Friend"`
	FriendData  account.Account   `json:"-" gorm:"foreignKey:Account"`
	FriendKey   account.PublicKey `json:"fr_key" gorm:"foreignKey:Account"`
	AccountKey  account.PublicKey `json:"acc_key" gorm:"foreignKey:Friend"`
}
