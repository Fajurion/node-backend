package manage

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type addRequest struct {
	Token string `json:"token"`

	Cluster         uint    `json:"cluster"` // Cluster ID
	IP              string  `json:"ip"`
	Port            uint    `json:"port"`
	PeformanceLevel float32 `json:"performance_level"`
}

func addNode(c *fiber.Ctx) error {

	// Parse body to add request
	var req addRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.FailedRequest(c, "invalid", err)
	}

	// Check if request is valid
	if req.Token == "" || req.Cluster == 0 || req.IP == "" || req.Port == 0 {
		return requests.FailedRequest(c, "invalid", nil)
	}

	// Check if IP is valid
	if len(req.IP) < 7 || len(req.IP) > 15 {
		return requests.FailedRequest(c, "invalid", nil)
	}

	// Check if port is valid
	if req.Port < 1 || req.Port > 65535 {
		return requests.FailedRequest(c, "invalid", nil)
	}

	// Check session and permission
	var session account.Session
	if !requests.CheckSessionPermission(c, req.Token, util.PermissionAdmin, &session) {
		return requests.FailedRequest(c, "no.permission", nil)
	}

	var cluster node.Cluster
	err := database.DBConn.First(&cluster, req.Cluster).Error

	if err != nil {
		return requests.FailedRequest(c, "invalid", err)
	}

	// Create node
	database.DBConn.Create(&node.Node{
		ClusterID:       req.Cluster,
		Token:           auth.GenerateToken(),
		Domain:          req.IP + ":" + strconv.FormatUint(uint64(req.Port), 10),
		PeformanceLevel: req.PeformanceLevel,
		Status:          1,
	})

	return c.SendString("Manage")
}
