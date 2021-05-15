package main

import (
	"projectGroup23/caching"
	"projectGroup23/database"
	"projectGroup23/discord"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initializes Firebase database
	database.InitFirebase()

	// Initializes AzureSQL database
	database.InitSQL()

	// Initializes the cache
	caching.InitCache()

	// Gets stored API response from last session
	database.GetStoredFromFirestore()

	// Initialize Discord Bot
	discord.InitDiscord()
}
