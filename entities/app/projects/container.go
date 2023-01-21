package projects

type Container struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Project uint   `json:"project"`
	Data    []byte `json:"data"`
}
