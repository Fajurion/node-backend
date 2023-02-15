package projects

type Member struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Account uint `json:"account"`
	Project uint `json:"project"`
	App     uint `json:"app"`
	Role    uint `json:"role"`

	ProjectData Project `json:"project_data" gorm:"foreignKey:Project"`
}
