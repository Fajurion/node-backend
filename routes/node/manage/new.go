package manage

import (
	"node-backend/database"
	"node-backend/entities/app"
	"node-backend/entities/node"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type newRequest struct {
	Token           string  `json:"token"`
	Cluster         uint    `json:"cluster"` // Cluster ID
	App             uint    `json:"app"`     // App ID
	Domain          string  `json:"domain"`
	PeformanceLevel float32 `json:"performance_level"`
}

func newNode(c *fiber.Ctx) error {

	// Parse body to add request
	var req newRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if token is valid
	var ct node.NodeCreation
	if err := database.DBConn.Take(&ct, req.Token).Error; err != nil {
		return requests.FailedRequest(c, "invalid", nil)
	}

	if req.Cluster == 0 || req.Domain == "" {
		return requests.FailedRequest(c, "invalid", nil)
	}

	if len(req.Domain) < 3 {
		return requests.FailedRequest(c, "invalid.domain", nil)
	}

	var cluster node.Cluster
	if err := database.DBConn.Take(&cluster, req.Cluster).Error; err != nil {
		return requests.FailedRequest(c, "invalid", nil)
	}

	var app app.App
	if err := database.DBConn.Take(&app, req.App).Error; err != nil {
		return requests.FailedRequest(c, "invalid", nil)
	}

	// Create node
	var created node.Node = node.Node{
		AppID:           req.App,
		ClusterID:       req.Cluster,
		Token:           auth.GenerateToken(300),
		Domain:          req.Domain,
		Load:            0,
		PeformanceLevel: req.PeformanceLevel,
		Status:          1,
	}
	database.DBConn.Create(&created)

	// Adopt new node
	var nodes []node.Node
	database.DBConn.Model(&node.Node{}).Where("cluster_id = ?", req.Cluster).Find(&nodes)

	for _, n := range nodes {
		if n.IsStarted() && n.ID != created.ID {
			n.SendAdoption(created.Domain, created.Token)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   created.Token,
		"cluster": cluster.Country,
		"app":     app.Name,
	})
}
