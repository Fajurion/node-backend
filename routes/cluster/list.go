package cluster

import (
	"node-backend/database"
	"node-backend/entities/node"

	"github.com/gofiber/fiber/v2"
)

func listClusters(c *fiber.Ctx) error {

	var clusters []node.Cluster
	database.DBConn.Find(&clusters)

	return c.JSON(clusters)
}
