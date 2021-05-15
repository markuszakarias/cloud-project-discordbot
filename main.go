package main

import (
	"projectGroup23/caching"
	"projectGroup23/database"
	"projectGroup23/discord"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database.InitFirebase()
	database.InitSQL()
	caching.InitCache()
	database.GetStoredFromFirestore()
	discord.RunDiscordBot()
}
