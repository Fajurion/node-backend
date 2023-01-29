package cluster

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type removeRequest struct {
	Token string `json:"token"`
	ID    uint   `json:"id"`
}

func removeCluster(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	if !util.Permission(c, util.PermissionAdmin) {
		return requests.InvalidRequest(c)
	}

	// Check if cluster exists
	var cluster node.Cluster
	err := database.DBConn.First(cluster, req.ID).Error

	if err == nil {
		return requests.FailedRequest(c, "cluster.exists", nil)
	}

	// Remove cluster
	err = database.DBConn.Delete(cluster).Error

	if err != nil {
		return requests.FailedRequest(c, "invalid", err)
	}

	return requests.SuccessfulRequest(c)
}
