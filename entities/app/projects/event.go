package projects

import (
	"time"
)

type Event struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Project uint      `json:"project"`
	Date    time.Time `json:"date"`
	Data    string    `json:"data"`
}
