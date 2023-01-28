package node

import "node-backend/entities/app"

type Node struct {
	ID uint `json:"id" gorm:"primaryKey"`

	AppID           uint    `json:"app"` // App ID
	ClusterID       uint    `json:"cluster"`
	Token           string  `json:"token"`
	Domain          string  `json:"domain" gorm:"unique"`
	Load            float64 `json:"load"`
	PeformanceLevel float32 `json:"performance_level"`

	// started: 0, stopped: 1, starting: 2, stopping: 3, error: 4
	Status uint `json:"status"`

	Cluster Cluster // This is an association (may still be broken)
	App     app.App // This is an association (may still be broken)
}

func (n *Node) IsStarted() bool {
	return n.Status == 0
}

func (n *Node) IsStopped() bool {
	return n.Status == 1
}

func (n *Node) IsStarting() bool {
	return n.Status == 2
}

func (n *Node) IsStopping() bool {
	return n.Status == 3
}

func (n *Node) HadError() bool {
	return n.Status == 4
}

func (n *Node) SendAdoption(domain string, token string) error {

	// Send adoption request

	return nil
}

func (n *Node) GetConnection(token string) error {

	// Get new token

	return nil
}