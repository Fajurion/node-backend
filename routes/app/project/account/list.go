package account

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/app/projects"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type listRequest struct {
	App    uint `json:"app"`
	LastID uint `json:"last_id"`
}

func listProjects(c *fiber.Ctx) error {

	// Parse request
	var req listRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Get projects
	data := util.GetData(c)
	acc := data["acc"].(uint)
	tk := util.GetToken(c)

	var requested account.Session
	if err := database.DBConn.Where("token = ?", tk).Take(&requested).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	var members []projects.Member
	if err := database.DBConn.Where("app = ? AND account = ? AND id > ?", requested.App, acc, req.LastID).Order("id DESC").Limit(10).Preload("ProjectData").Find(&members).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	var projects []projects.ProjectEntity
	for _, member := range members {
		projects = append(projects, member.ProjectData.ToEntity())
	}

	// Return response
	return c.JSON(fiber.Map{
		"success":  true,
		"projects": projects,
	})
}
