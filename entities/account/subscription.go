package account

import "time"

type Subscription struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Account   string    `json:"account"`
	Rank      uint      `json:"rank"`
	End       time.Time `json:"end"`
	Price     float32   `json:"price"`
	Duration  uint      `json:"duration"`
	Renewable bool      `json:"renewable"`
}
