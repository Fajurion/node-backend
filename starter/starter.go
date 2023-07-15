package starter

import (
	"bufio"
	"log"
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/routes"
	"node-backend/util"
	"node-backend/util/nodes"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func Startup() {

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	util.JWT_SECRET = os.Getenv("JWT_SECRET")

	// Connect to the database
	database.Connect()

	app.Use(cors.New())
	app.Use(logger.New())

	// Handle routing
	app.Route("/", routes.Router)

	// Ask user for test mode
	testMode()

	// Listen on port 3000
	err = app.Listen(os.Getenv("LISTEN"))

	log.Println(err.Error())
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

	token, _ := util.Token("123", "123", 100, time.Now().Add(time.Hour*24))

	log.Println("Test token: " + token)

	var foundNodes []node.Node
	database.DBConn.Find(&foundNodes)

	for _, n := range foundNodes {
		if n.Status == node.StatusStarted {
			log.Println("Stopping node", n.Domain)

			nodes.TurnOff(&n, node.StatusStopped)
		}
	}
}
