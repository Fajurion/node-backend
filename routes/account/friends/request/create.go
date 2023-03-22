package request

import (
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/entities/node"
	"node-backend/util/nodes"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type addFriendRequest struct {
	Node      uint   `json:"id"`
	NodeToken string `json:"token"`
	Session   uint   `json:"session"`
	Username  string `json:"username"`
	Tag       string `json:"tag"`
}

type addFriendResponse struct {
	Success bool            `json:"success"`
	Action  string          `json:"action"`
	Friend  uint            `json:"friend"`
	Node    node.NodeEntity `json:"node"`
	Key     string          `json:"key"`
}

// Route: /account/friends/request/create
func createRequest(c *fiber.Ctx) error {

	var req addFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return requests.InvalidRequest(c)
	}

	_, err := nodes.Node(req.Node, req.NodeToken)
	if err != nil {
		return requests.InvalidRequest(c)
	}

	var session account.Session
	if !requests.GetSession(req.Session, &session) {
		return requests.InvalidRequest(c)
	}

	log.Println(session.Account)

	var friend account.Account
	if err := database.DBConn.Where(&account.Account{Username: req.Username, Tag: req.Tag}).Take(&friend).Error; err != nil {
		return requests.InvalidRequest(c)
	}

	if friend.ID == session.Account {
		return requests.FailedRequest(c, "cannot.add.self", nil)
	}

	// Check if the friend is already a friend
	var friendCheck properties.Friend
	if err := database.DBConn.Where(&properties.Friend{Account: friend.ID, Friend: session.Account}).Take(&friendCheck).Error; err == nil && !friendCheck.Request {
		return requests.FailedRequest(c, "already.friends", nil)
	}

	if err := database.DBConn.Where(&properties.Friend{Account: session.Account, Friend: friend.ID}).Take(&properties.Friend{}).Error; err == nil {
		return requests.FailedRequest(c, "already.requested", nil)
	}

	var friendSession account.Session
	database.DBConn.Where(&account.Session{Account: friend.ID}).Not("node = ?", 0).Take(&friendSession) // Doesn't matter if the session is connected or not

	if friendCheck.Request {

		// Accept friend request
		friendCheck.Request = false
		if err := database.DBConn.Omit("Friend").Save(&friendCheck).Error; err != nil {
			return requests.FailedRequest(c, "server.error", err)
		}

		database.DBConn.Create(&properties.Friend{
			Account: session.Account,
			Friend:  friend.ID,
			Request: false,
		})

		return ExecuteAction(c, "accept", friend.ID, friendSession)
	}

	// Send friend request
	if err := database.DBConn.Create(&properties.Friend{
		Account: session.Account,
		Friend:  friend.ID,
		Request: true,
	}).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Send notification to friend
	return ExecuteAction(c, "send", friend.ID, friendSession)
}

// ExecuteAction returns the action to the node
func ExecuteAction(c *fiber.Ctx, action string, friend uint, session account.Session) error {

	if session.Token == "" {
		return c.JSON(addFriendResponse{
			Success: true,
			Action:  action,
			Friend:  friend,
			Node:    node.NodeEntity{},
			Key:     "",
		})
	}

	var nodeInfo node.Node
	if err := database.DBConn.Where(&node.Node{ID: session.Node}).Take(&nodeInfo).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(addFriendResponse{
		Success: true,
		Action:  action,
		Friend:  friend,
		Node:    nodeInfo.ToEntity(),
	})

}
