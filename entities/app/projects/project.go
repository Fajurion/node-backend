package projects

type Project struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Creator uint   `json:"creator"`
	App     uint   `json:"app"`
	Data    string `json:"data"`

	Events     []Event     `json:"events" gorm:"foreignKey:Project"`
	Containers []Container `json:"containers" gorm:"foreignKey:Project"`
	Members    []Member    `json:"members" gorm:"foreignKey:Project"`
}

func (p Project) ToEntity() ProjectEntity {
	return ProjectEntity{
		ID:   p.ID,
		Data: p.Data,
	}
}

type ProjectEntity struct {
	ID   uint   `json:"id"`
	Data string `json:"data"`
}
