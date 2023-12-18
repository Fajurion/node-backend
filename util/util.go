package util

import (
	"crypto/rsa"
	"io"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

var Testing = false
var LogErrors = true

const PermissionUseServices = 10
const PermissionServicesUnlimited = 20
const PermissionViewSettings = 50
const PermissionManageNodes = 55
const PermissionChangeSettings = 60
const PermissionAdmin = 100

var JWT_SECRET = "hi"

func PostRequest(url string, body map[string]interface{}) (map[string]interface{}, error) {

	req, _ := sonic.Marshal(body)

	reader := strings.NewReader(string(req))

	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		return nil, err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)

	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = sonic.Unmarshal([]byte(buf.String()), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Parse encrypted json
func BodyParser(c *fiber.Ctx, data interface{}) error {
	return sonic.Unmarshal(c.Locals("body").([]byte), data)
}

// Return encrypted json
func ReturnJSON(c *fiber.Ctx, data interface{}) error {

	encoded, err := sonic.Marshal(data)
	if err != nil {
		return FailedRequest(c, ErrorServer, err)
	}

	encrypted, err := EncryptRSA(c.Locals("pub").(*rsa.PublicKey), encoded)
	if err != nil {
		return FailedRequest(c, ErrorServer, err)
	}

	return c.Send(encrypted)
}
