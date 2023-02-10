package util

import (
	"io"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
)

const PermissionUseServices = 10
const PermissionServicesUnlimited = 20
const PermissionViewSettings = 50
const PermissionManageNodes = 55
const PermissionChangeSettings = 60
const PermissionAdmin = 100

const JWT_SECRET = "prJ9Qe6TQhLNn5PsEGkA7qAY9oDr7K349gKrrnHTHJKaYJkc6hpzJDGqrRFPTSLDs4oa"
const NODE_PW = "RDM1LpYetf7Yw3ZTv40hjdeRfbiIkcyGbSDsm4xyN3zIBwCMKDRi9uPxN8SwjbTfyYtrQ9hyyOlYkCXn4mQVlllB73ZsxW6OL8wXyFtUBGi3wdq8mkQNYpUiRzkSCIdfDpOiUeb7DxqVZaKHAxhDelxZ4BquC6VKIewmQy67s7UwONPuFacqNQd8G58zv3YTrEtoWSzngMcieGHL9snIgsNk5SHyHHy00VIAz9z6aXANHgg3yFC48yQTKzvHBKqd5sbpAkI5l86Yqrk2P7uV75ZQBr5l62mHVxJyEnvV41Uf"

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
