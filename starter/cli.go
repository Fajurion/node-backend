package starter

import (
	"bufio"
	"fmt"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/app"
	"node-backend/entities/node"
	"node-backend/util/auth"
	"os"
	"strconv"
	"strings"
	"time"
)

func listenForCommands() {
	for {
		fmt.Print("node-backend > ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "exit":
			os.Exit(0)
		case "create-default":

			if database.DBConn.Where("name = ?", "Default").Take(&account.Rank{}).RowsAffected > 0 {
				fmt.Println("Default stuff already exists")
				continue
			}

			// Create default ranks
			database.DBConn.Create(&account.Rank{
				Name:  "Default",
				Level: 20,
			})
			database.DBConn.Create(&account.Rank{
				Name:  "Admin",
				Level: 100,
			})

			// Create default cluster
			database.DBConn.Create(&node.Cluster{
				Name:    "Vaterland",
				Country: "DE",
			})

			fmt.Println("Created default ranks and cluster")

		case "create-app":

			fmt.Print("App name: ")
			appName, _ := reader.ReadString('\n')
			appName = strings.TrimSpace(appName)

			fmt.Print("App description: ")
			appDescription, _ := reader.ReadString('\n')
			appDescription = strings.TrimSpace(appDescription)

			fmt.Print("App version: ")
			appVersion, _ := reader.ReadString('\n')
			appVersion = strings.TrimSpace(appVersion)

			fmt.Print("App access level: ")
			appAccessLevel, _ := reader.ReadString('\n')
			appAccessLevel = strings.TrimSpace(appAccessLevel)

			// Create app
			accessLevel, err := strconv.Atoi(appAccessLevel)
			if err != nil {
				fmt.Println("Invalid access level")
				continue
			}

			app := &app.App{
				Name:        appName,
				Description: appDescription,
				Version:     appVersion,
				AccessLevel: uint(accessLevel),
			}
			database.DBConn.Create(&app)

			fmt.Println("Created app with ID", app.ID)

		case "create-node":

			// Generate new node token
			tk := auth.GenerateToken(100)

			// Save
			if err := database.DBConn.Create(&node.NodeCreation{
				Token: tk,
				Date:  time.Now(),
			}).Error; err != nil {
				fmt.Println("Failed to create node token")
				continue
			}

			fmt.Println("Created node token", tk)

		case "help":
			fmt.Println("exit - Exit the CLI")
			fmt.Println("create-default - Create default ranks and cluster")
			fmt.Println("create-app - Create a new app")
			fmt.Println("create-node - Get a node token (rest of setup in the CLI of the node)")

		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
