package node

import (
	"io"
	"net/http"
	"node-backend/entities/app"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

type Node struct {
	ID uint `json:"id" gorm:"primaryKey"`

	AppID           uint    `json:"app"` // App ID
	ClusterID       uint    `json:"cluster"`
	Token           string  `json:"token"`
	Domain          string  `json:"domain" gorm:"unique"`
	Load            float64 `json:"load"`
	PeformanceLevel float32 `json:"performance_level"`

	// 0 - Started, 1 - Stopped, 2 - Error
	Status uint `json:"status"`

	Cluster Cluster // This is an association (may still be broken)
	App     app.App // This is an association (may still be broken)
}

func (n *Node) SendAdoption(domain string, token string) error {

	// Send adoption request

	return nil
}

func (n *Node) GetConnection(token string, user uint) (string, error) {

	req, _ := sonic.Marshal(fiber.Map{
		"node_token": n.Token,
		"session":    token,
		"user_id":    user,
	})

	reader := strings.NewReader(string(req))

	res, err := http.Post("http://"+n.Domain+"/auth/initialize", "application/json", reader)
	if err != nil {
		return "", err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)

	if err != nil {
		return "", err
	}

	var data map[string]string
	err = sonic.Unmarshal([]byte(buf.String()), &data)
	if err != nil {
		return "", err
	}

	return data["token"], nil
}

const StatusStarted = 0
const StatusStopped = 1
const StatusError = 2
