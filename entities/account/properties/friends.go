package properties

import "time"

type Friend struct {
	Account uint      `json:"account" gorm:"primaryKey"`
	Friend  uint      `json:"friend"`
	Date    time.Time `json:"date"`
	Request bool      `json:"request"`
}
