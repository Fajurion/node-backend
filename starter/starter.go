package starter

import (
	"bufio"
	"fmt"
	"log"
	"net/smtp"
	"node-backend/database"
	routes_v1 "node-backend/routes/v1"
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
		JSONEncoder:       sonic.Marshal,
		JSONDecoder:       sonic.Unmarshal,
		StreamRequestBody: true, // TODO: Proper request body protection (Make only certain endpoints accept streams)
	})

	util.TestAES()

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	util.JWT_SECRET = os.Getenv("JWT_SECRET")

	/* to test email
	err = testMail()
	if err != nil {
		panic(err)
	}
	*/

	// Connect to the databases
	database.Connect()

	app.Use(cors.New())
	app.Use(logger.New())

	// Handle routing
	app.Route("/v1", routes_v1.Router)

	// Ask user for test mode
	testMode()

	// Listen on port 3000
	if os.Getenv("CLI") == "true" {
		go func() {
			err = app.Listen(os.Getenv("LISTEN"))

			log.Println(err.Error())
		}()

		// Listen for commands
		listenForCommands()
	} else {
		err = app.Listen(os.Getenv("LISTEN"))

		log.Println(err.Error())
	}
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

func testMail() error {

	// Generate message
	subject := "Chat app test email"
	body := "Hello, this is a test email sent through Amazon SES by the Chat app server."

	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r%s", subject, body))

	// Authenticate using the provided credentials
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PW"), os.Getenv("SMTP_SERVER"))

	// Send the email
	err := smtp.SendMail(
		os.Getenv("SMTP_SERVER")+":"+os.Getenv("SMTP_PORT"),
		auth,
		os.Getenv("SMTP_FROM"),
		[]string{"test@email.com"},
		msg,
	)
	return err
}
