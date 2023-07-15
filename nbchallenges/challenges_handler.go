package nb_challenges

import (
	"node-backend/util/requests"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type generateRequest struct {
	Token string `json:"challenge"`
}

var challengeMap = map[string]string{}

func Generate(c *fiber.Ctx) error {

	var req generateRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	// Generate challenge
	tk, result, js, attach := GenerationFunction()
	js = strings.ReplaceAll(js, "\n", "")
	js = strings.ReplaceAll(js, "\r", "")

	challengeMap[tk] = result

	return c.JSON(fiber.Map{
		"success": true,
		"tk":      tk,
		"js":      js,
		"attach":  attach,
	})
}

type solveRequest struct {
	Token  string `json:"challenge"`
	Result string `json:"result"`
}

func Solve(c *fiber.Ctx) error {
	var req solveRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	if req.Result != challengeMap[req.Token] {
		return c.JSON(fiber.Map{
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
