package account

type Authentication struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Account string `json:"account"`
	Type    uint   `json:"type"`
	Secret  string `json:"secret"`
}
