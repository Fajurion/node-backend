package app

import "node-backend/entities/node"

type AppNode struct {
	ID uint `json:"id" gorm:"primaryKey"`

	NodeID           uint    `json:"node"`
	App              uint    `json:"app"`
	Load             uint    `json:"load"`
	PerformanceLevel float64 `json:"performance_level"`
	Cluster          uint    `json:"cluster"`

	Node node.Node `gorm:"foreignKey:Node"`
}
