package properties

type Friend struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Account   uint   `json:"account" gorm:"primaryKey"`
	Friend    uint   `json:"friend"`
	Updated   int64  `gorm:"autoUpdateTime:milli"`
	Request   bool   `json:"request"`
	Signature string `json:"signature"`
}
