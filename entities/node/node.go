package node

import (
	"errors"
	"node-backend/entities/account"
	"node-backend/entities/app"
	"node-backend/util"

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

	// 1 - Started, 2 - Stopped, 3 - Error
	Status uint `json:"status"`

	Cluster Cluster `gorm:"foreignKey:ClusterID"`
	App     app.App `gorm:"foreignKey:AppID"`
}

func (n *Node) ToEntity() NodeEntity {
	return NodeEntity{
		ID:     n.ID,
		Token:  n.Token,
		App:    n.AppID,
		Domain: n.Domain,
	}
}

func (n *Node) SendPing(node Node) error {

	_, err := util.PostRequest("http://"+n.Domain+"/ping", fiber.Map{})

	return err
}

func (n *Node) GetConnection(acc account.Account, session uint, sessionIds []uint) (string, bool, error) {

	// Check for key
	if acc.Key.Key == "" {
		return "", true, errors.New("no.key")
	}

	// Send request to node
	res, err := util.PostRequest("http://"+n.Domain+"/auth/initialize", fiber.Map{
		"node_token":  n.Token,
		"session":     session,
		"user_id":     acc.ID,
		"username":    acc.Username,
		"tag":         acc.Tag,
		"session_ids": sessionIds,
	})

	if err != nil {
		return "", false, err
	}

	if res["load"] != nil {
		n.Load = res["load"].(float64)
	}

	if !res["success"].(bool) {
		return "", true, errors.New(res["message"].(string))
	}

	return res["token"].(string), true, nil

	/* PREVIOUS CODE
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

	var data map[string]interface{}
	err = sonic.Unmarshal([]byte(buf.String()), &data)
	if err != nil {
		return "", err
	}

	n.Load = data["load"].(float64)

	return data["token"].(string), nil
	*/
}

const StatusStarted = 1
const StatusStopped = 2
const StatusError = 3

type NodeEntity struct {
	ID     uint   `json:"id"`
	Token  string `json:"token"`
	App    uint   `json:"app"`
	Domain string `json:"domain"`
}
