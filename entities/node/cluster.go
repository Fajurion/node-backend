package node

type Cluster struct {
	ID string `json:"id" gorm:"primaryKey"`

	Name    string `json:"name"`
	Country string `json:"country"`
}
