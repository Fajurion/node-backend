package node

import (
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

	// 0 - Started, 1 - Stopped, 2 - Error
	Status uint `json:"status"`

	Cluster Cluster // This is an association (may still be broken)
	App     app.App // This is an association (may still be broken)
}

func (n *Node) ToEntity() NodeEntity {
	return NodeEntity{
		ID:     n.ID,
		Token:  n.Token,
		App:    n.AppID,
		Domain: n.Domain,
	}
}

func (n *Node) SendAdoption(node Node) error {

	_, err := util.PostRequest("http://"+n.Domain+"/adoption/initialize", fiber.Map{
		"token": n.Token,
		"node":  node.ToEntity(),
	})

	if err != nil {
		return err
	}

	return nil
}

func (n *Node) GetConnection(token string, user uint) (string, error) {

	res, err := util.PostRequest("http://"+n.Domain+"/auth/initialize", fiber.Map{
		"node_token": n.Token,
		"session":    token,
		"user_id":    user,
	})

	if err != nil {
		return "", err
	}

	if res["load"] != nil {
		n.Load = res["load"].(float64)
	}

	return res["token"].(string), nil

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

const StatusStarted = 0
const StatusStopped = 1
const StatusError = 2

type NodeEntity struct {
	ID     uint   `json:"id"`
	Token  string `json:"token"`
	App    uint   `json:"app"`
	Domain string `json:"domain"`
}
