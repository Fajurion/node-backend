package node

import (
	"errors"
	"node-backend/util"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Node struct {
	ID uint `json:"id" gorm:"primaryKey"`

	AppID           uint    `json:"app"` // App ID
	ClusterID       uint    `json:"cluster"`
	Token           string  `json:"token"`
	Domain          string  `json:"domain"`
	Load            float64 `json:"load"`
	PeformanceLevel float32 `json:"performance_level"`

	// 1 - Started, 2 - Stopped, 3 - Error
	Status uint `json:"status"`
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

	_, err := util.PostRequestNoTC("http://"+n.Domain+"/ping", map[string]interface{}{})

	return err
}

// Sender for GetConnection
const SenderUser = 0
const SenderNode = 1

func (n *Node) GetConnection(accId string, session string, sessionIds []string, sender int) (string, bool, error) {

	if sender != SenderUser && sender != SenderNode {
		return "", false, errors.New("invalid.sender")
	}

	// Get public key of node
	res, err := util.PostRequestNoTC(util.NodeProtocol+n.Domain+"/pub", fiber.Map{})
	if err != nil {
		return "", false, err
	}

	// Unpackage the public key
	publicKey, err := util.UnpackageRSAPublicKey(res["pub"].(string))
	if err != nil {
		return "", false, err
	}

	// Send request to node
	res, err = util.PostRequest(publicKey, util.NodeProtocol+n.Domain+"/auth/initialize", fiber.Map{
		"sender":     sender,
		"node_token": n.Token,
		"session":    session,
		"account":    accId,
		"end":        time.Now().UnixMilli(),
	})
	if err != nil {
		return "", false, err
	}

	// Set the new load (will later be updated in database)
	if res["load"] != nil {
		n.Load = res["load"].(float64)
	}

	// Check if request was successful
	if !res["success"].(bool) {
		return "", true, errors.New(res["message"].(string))
	}

	// Return connection token
	return res["token"].(string), true, nil
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
