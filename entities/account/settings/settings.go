package settings

type AccountSetting struct {
	Account uint   `json:"account" gorm:"primaryKey"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}
