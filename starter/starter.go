package starter

import (
	"bufio"
	"log"
	"node-backend/database"
	"node-backend/routes"
	"node-backend/util"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func Startup() {

	log.SetOutput(os.Stdout)

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

	if os.Getenv("TESTING") != "" {
		util.Testing = os.Getenv("TESTING") == "true"
		if util.Testing {
			log.Println("Test mode enabled (Read from .env).")
		}
	} else {
		log.Println("Do you want to continue in test mode? (y/n)")

		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		util.Testing = scanner.Text() == "y"
	}

	if !util.Testing {
		return
	}

	log.Println("Test mode enabled.")

	token, _ := util.Token("123", "123", 100, time.Now().Add(time.Hour*24))

	log.Println("Test token: " + token)

	/* not need for now
	var foundNodes []node.Node
	database.DBConn.Find(&foundNodes)

	for _, n := range foundNodes {
		if n.Status == node.StatusStarted {
			log.Println("Stopping node", n.Domain)

			nodes.TurnOff(&n, node.StatusStopped)
		}
	}
	*/
}
