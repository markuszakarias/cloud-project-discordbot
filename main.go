package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"projectGroup23/caching"
	"projectGroup23/database"
	"projectGroup23/discordpkg/discordutils"
	"projectGroup23/webhook"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/bwmarrin/discordgo"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

var server = envVar("DB_SERVER")
var port = envVar("DB_PORT")
var user = envVar("DB_USER")
var password = envVar("DB_PASSWORD")
var db = envVar("DB")

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, db)

	// Initializes Firebase database
	database.InitFirebase()

	// Initializes BigCache cache
	err := caching.AddCacheModule("cache")
	if err != nil {
		fmt.Println("Error with initialize cache: " + err.Error())
	}

	// Initializes BigCache cache
	err = caching.AddCacheModule("weather")
	if err != nil {
		fmt.Println("Error with initialize cache: " + err.Error())
	}

	// Gets stored API response from last session
	database.GetStoredFromFirestore()

	// Initializes Discord bot with token
	token := envVar("DC_TOKEN")
	s, err := discordgo.New("Bot " + token)

	if err = s.Open(); err != nil {
		panic(err)
	}

	go webhook.WebhookRoutine(s)
	s.AddHandler(messageCreate)
	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Create connection pool
	database.Db, database.Err = sql.Open("sqlserver", connString)
	if database.Err != nil {
		log.Fatal("Error creating connection pool: ", database.Err.Error())
	}
	ctx := context.Background()
	database.Err = database.Db.PingContext(ctx)
	if database.Err != nil {
		log.Fatal(database.Err.Error())
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	s.Close()
}

// envVar reads from .env file, sets environment variables and returns value based on key
func envVar(key string) string {
	// Load .env file
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println(err)
	}
	return os.Getenv(key)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content[:1] != "!" {
		return
	}
	var err error

	switch {
	case len(m.Content) >= 5 && m.Content[:5] == "!help":
		discordutils.SendHelpMessage(s, m)
	case m.Content == "!steamdeals":
		dur, _ := time.ParseDuration("20s")
		err = caching.CacheDeals(m.Content, dur)
		discordutils.SendSteamMessage(s, m)
	case len(m.Content) >= 8 && m.Content[:8] == "!weather":
		err = discordutils.SendWeatherMessage(s, m)
	case m.Content == "!mealplan":
		dur, _ := time.ParseDuration("20s")
		err = caching.CacheMeals(dur)
		discordutils.SendMealplanMessage(s, m)
	case len(m.Content) >= 11 && m.Content[:11] == "!newsletter":
		/* dur, _ := time.ParseDuration("20s")
		err = caching.CacheNews(dur) */
		err = discordutils.SendNewsletterMessage(s, m)
	case len(m.Content) >= 5 && m.Content[:5] == "!todo":
		discordutils.SendTodoMessage(s, m)
	case len(m.Content) >= 14 && m.Content[:14] == "!notifyweather":
		err = discordutils.NotifyWeather(s, m)
	case len(m.Content) >= 5 && m.Content[:5] == "!joke":
		err = discordutils.SendJokeMessage(s, m)

	default:
		s.ChannelMessageSend(m.ChannelID, "Unable to recognize command, try !help (not implemented) if you need a reminder!")
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
	}
}
