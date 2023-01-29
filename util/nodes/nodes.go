package nodes

import (
	"node-backend/database"
	"node-backend/entities/node"
)

func Node(token string) (node.Node, error) {

	// Check if token is valid
	var found node.Node
	if err := database.DBConn.Where(&node.Node{Token: token}).Take(&found).Error; err != nil {
		return node.Node{}, err
	}

	return found, nil
}
