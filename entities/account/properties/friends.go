package properties

import "node-backend/entities/account"

type Friend struct {
	Account uint  `json:"account" gorm:"primaryKey"`
	Friend  uint  `json:"friend"`
	Updated int64 `gorm:"autoUpdateTime:milli"`
	Request bool  `json:"request"`

	AccountData account.Account `json:"-" gorm:"foreignKey:Friend"`
	FriendData  account.Account `json:"-" gorm:"foreignKey:Account"`
}
