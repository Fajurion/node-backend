package properties

type Friend struct {
	Account uint `json:"account" gorm:"primaryKey"`
	Friend  uint `json:"friend"`
	Request bool `json:"request"`
}
