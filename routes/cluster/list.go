package cluster

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

func listClusters(c *fiber.Ctx) error {

	if !util.Permission(c, util.PermissionUseServices) {
		return requests.InvalidRequest(c)
	}

	var clusters []node.Cluster
	database.DBConn.Find(&clusters)

	return c.JSON(clusters)
}
