package main

import (
	"projectGroup23/caching"
	"projectGroup23/database"
	"projectGroup23/discord"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	database.InitFirebase()

	database.InitSQL()

	caching.InitCache()
}

func main() {
	database.GetStoredFromFirestore()
	discord.RunDiscordBot()
}
