package app

import "node-backend/entities/app/projects"

type App struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	AccessLevel uint   `json:"access_level"`

	Projects []projects.Project `json:"projects" gorm:"foreignKey:App"`
}
