package cluster

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type addRequest struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

func addCluster(c *fiber.Ctx) error {

	// Parse request
	var req addRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	if !util.Permission(c, util.PermissionAdmin) {
		return requests.InvalidRequest(c)
	}

	// Check if cluster name is valid
	if len(req.Name) < 3 {
		return requests.FailedRequest(c, "invalid.name", nil)
	}

	// Check if country is valid
	if len(req.Country) != 2 {
		return requests.FailedRequest(c, "invalid.country", nil)
	}

	// Check if cluster already exists
	err := database.DBConn.Create(&node.Cluster{
		Name:    req.Name,
		Country: req.Country,
	}).Error

	if err != nil {
		return requests.FailedRequest(c, "cluster.exists", err)
	}

	return requests.SuccessfulRequest(c)
}
