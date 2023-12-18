package nb_challenges

import (
	"node-backend/util"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type generateRequest struct {
	Token string `json:"challenge"`
}

var challengeMap = map[string]string{}

func Generate(c *fiber.Ctx) error {

	var req generateRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Generate challenge
	tk, result, js, attach := GenerationFunction()
	js = strings.ReplaceAll(js, "\n", "")
	js = strings.ReplaceAll(js, "\r", "")

	challengeMap[tk] = result

	return util.ReturnJSON(c, fiber.Map{
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	if req.Result != challengeMap[req.Token] {
		return util.ReturnJSON(c, fiber.Map{
			"success": false,
		})
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
	})
}
