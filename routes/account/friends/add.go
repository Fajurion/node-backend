package friends

import (
	"log"
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type addFriendRequest struct {
	Payload string `json:"payload"` // Encrypted payload
}

// Route: /account/friends/add
func addFriend(c *fiber.Ctx) error {

	log.Println(string(c.Body()))

	// Parse request
	var req addFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Check if the account has too many friends
	accId := util.GetAcc(c)
	var friendCount int64
	if err := database.DBConn.Model(&properties.Friendship{}).Where("account = ?", accId).Count(&friendCount).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	if friendCount >= MaximumFriends {
		return requests.FailedRequest(c, "limit.reached", nil)
	}

	// Create friendship
	friendship := properties.Friendship{
		ID:      auth.GenerateToken(12),
		Account: accId,
		Payload: req.Payload,
	}
	if err := database.DBConn.Create(&friendship).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"id":      friendship.ID,
	})
}
