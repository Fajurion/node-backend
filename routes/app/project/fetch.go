package project

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/app/projects"
	"node-backend/entities/node"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type fetchRequest struct {
	Token   string `json:"token"`
	Project uint   `json:"project"`
}

func fetch(c *fiber.Ctx) error {

	// Parse request
	var req fetchRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	var node node.Node
	if err := database.DBConn.Where("token = ?", req.Token).Take(&node).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	var project projects.Project
	if err := database.DBConn.Where("id = ?", req.Project).Preload("Members").Take(&project).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	if project.App != node.AppID {
		return requests.InvalidRequest(c)
	}

	var members []uint
	for _, member := range project.Members {
		members = append(members, member.ID)
	}

	var sessions []account.Session
	if err := database.DBConn.Model(&account.Session{}).Where("account IN ?", members).Find(&sessions).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	var memberNodes map[uint]uint = make(map[uint]uint)
	for _, session := range sessions {
		memberNodes[session.Account] = session.Node
	}

	return c.JSON(fiber.Map{
		"success": true,
		"project": project,
		"members": memberNodes,
	})
}
