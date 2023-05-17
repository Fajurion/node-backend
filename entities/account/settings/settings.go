package settings

type AccountSetting struct {
	Account string `json:"account" gorm:"primaryKey"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}
