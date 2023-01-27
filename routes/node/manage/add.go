package manage

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type addRequest struct {
	Cluster         uint    `json:"cluster"` // Cluster ID
	IP              string  `json:"ip"`
	Port            uint    `json:"port"`
	PeformanceLevel float32 `json:"performance_level"`
}

func addNode(c *fiber.Ctx) error {

	// Parse body to add request
	var req addRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	if !util.Permission(c, util.PermissionAdmin) {
		return requests.InvalidRequest(c)
	}

	if req.Cluster == 0 || req.IP == "" || req.Port == 0 {
		return requests.FailedRequest(c, "invalid", nil)
	}

	if len(req.IP) < 3 {
		return requests.FailedRequest(c, "invalid.ip", nil)
	}

	if req.Port < 1 || req.Port > 65535 {
		return requests.FailedRequest(c, "invalid.port", nil)
	}

	var cluster node.Cluster
	err := database.DBConn.Take(&cluster, req.Cluster).Error

	if err != nil {
		return requests.FailedRequest(c, "invalid", err)
	}

	// Create node
	database.DBConn.Create(&node.Node{
		ClusterID:       req.Cluster,
		Token:           auth.GenerateToken(300),
		Domain:          req.IP + ":" + strconv.FormatUint(uint64(req.Port), 10),
		PeformanceLevel: req.PeformanceLevel,
		Status:          1,
	})

	// Adopt new node
	var nodes []node.Node
	database.DBConn.Model(&node.Node{}).Where("cluster_id = ?", req.Cluster).Find(&nodes)

	for _, n := range nodes {
		if n.IsStarted() {
			n.SendAdoption()
		}
	}

	return c.SendString("Manage")
}
