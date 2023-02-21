package properties

type Friend struct {
	Account uint  `json:"account" gorm:"primaryKey"`
	Friend  uint  `json:"friend"`
	Updated int64 `gorm:"autoUpdateTime:milli"`
	Request bool  `json:"request"`
}
