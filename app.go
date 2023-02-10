package main

import (
	"bufio"
	"log"
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/routes"
	"node-backend/util"
	"node-backend/util/auth"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Connect to the database
	database.Connect()

	app.Use(cors.New())
	app.Use(logger.New())

	// Handle routing
	app.Route("/", routes.Router)

	// Comment this out when in production
	testMode()

	// Listen on port 3000
	app.Listen("127.0.0.1:3000")
}

func testMode() {

	log.Print("\n TEST MODE ENABLED \n")
	log.Println("Do you want to continue in test mode? (y/n)")

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	input := scanner.Text()

	if input != "y" {
		return
	}

	token, _ := util.Token(auth.GenerateToken(300), time.Now().Add(time.Hour*24), map[string]interface{}{
		"acc": 123,
		"lvl": 100,
	})

	log.Println("Test token: " + token)

	var nodes []node.Node
	database.DBConn.Find(&nodes)

	for _, n := range nodes {
		if n.Status == node.StatusStarted {
			log.Println("Stopping node", n.Domain)
			n.Status = node.StatusStopped

			database.DBConn.Save(&n)
		}
	}
}
