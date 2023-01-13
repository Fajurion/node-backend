package account

type Rank struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Name  string `json:"name"`
	Level uint   `json:"level"`
}
