package starter

import (
	"bufio"
	"fmt"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/entities/account/properties"
	"node-backend/entities/app"
	"node-backend/entities/node"
	"node-backend/util"
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

		case "delete-data":

			fmt.Print("Account E-Mail: ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)

			// Delete all data
			var acc account.Account
			if err := database.DBConn.Where("email = ?", email).Take(&acc).Error; err != nil {
				fmt.Println("Failed to find account")
				continue
			}

			database.DBConn.Where("account = ?", acc.ID).Delete(&account.Session{})
			database.DBConn.Where("id = ?", acc.ID).Delete(&account.ProfileKey{})
			database.DBConn.Where("id = ?", acc.ID).Delete(&account.StoredActionKey{})
			database.DBConn.Where("id = ?", acc.ID).Delete(&account.PublicKey{})
			database.DBConn.Where("account = ?", acc.ID).Delete(&properties.AStoredAction{})
			database.DBConn.Where("account = ?", acc.ID).Delete(&properties.StoredAction{})
			database.DBConn.Where("account = ?", acc.ID).Delete(&properties.Friendship{})
			database.DBConn.Where("account = ?", acc.ID).Delete(&properties.VaultEntry{})

		case "account-token":

			fmt.Print("Account E-Mail: ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)

			// Get account
			var acc account.Account
			if err := database.DBConn.Where("email = ?", email).Preload("Rank").Take(&acc).Error; err != nil {
				fmt.Println("Failed to find account")
				continue
			}

			// Generate new token
			util.Token("test", acc.ID, acc.Rank.Level, time.Now().Add(time.Hour*24*365))

		case "help":
			fmt.Println("exit - Exit the CLI")
			fmt.Println("create-default - Create default ranks and cluster")
			fmt.Println("create-app - Create a new app")
			fmt.Println("create-node - Get a node token (rest of setup in the CLI of the node)")
			fmt.Println("delete-data - Delete the data to restart the setup process on an account")
			fmt.Println("account-token - Generate a JWT token for an account")

		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
