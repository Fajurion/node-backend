package util

import (
	"io"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
)

var Testing = false

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
