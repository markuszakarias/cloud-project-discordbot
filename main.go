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

	//go database.WebhookRoutine(s)
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
	err := godotenv.Load(".env")
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

	switch {
	case m.Content[:5] == "!help":
		discordutils.SendHelpMessage(s, m)
	case m.Content == "!steamdeals":
		dur, _ := time.ParseDuration("20s")
		err := caching.CacheDeals(m.Content, dur)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		}
		discordutils.SendSteamMessage(s, m)
	case m.Content[:8] == "!weather":
		discordutils.SendWeatherMessage(s, m)
	case m.Content == "!mealplan":
		dur, _ := time.ParseDuration("20s")
		err := caching.CacheMeals(dur)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		}
		discordutils.SendMealplanMessage(s, m)
	case m.Content == "!newsletter":
		dur, _ := time.ParseDuration("20s")
		err := caching.CacheNews(dur)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		}
		discordutils.SendNewsletterMessage(s, m)
	case m.Content[:5] == "!todo":
		discordutils.SendTodoMessage(s, m)
	case m.Content[:14] == "!notifyweather":
		err := discordutils.NotifyWeather(s, m)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		}
	default:
		s.ChannelMessageSend(m.ChannelID, "Unable to recognize command, try !help (not implemented) if you need a reminder!")
	}

	// TODO Incorporate notifications into switch
	/*
		if m.Content == "!notifyweather remove" {
			err := database.DeleteWebhook(m.Author.ID)
			if err != nil {
				return
			}
			s.ChannelMessageSend(m.ChannelID, "cloud notification removed!")
		}

		if m.Content[:21] == "!notifyweather cloud " {
			percent, err := strconv.Atoi(m.Content[21:])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "cloud percentage needs to be a number")
				return
			}
			if percent < 0 || 100 < percent {
				s.ChannelMessageSend(m.ChannelID, "cloud percentage needs to beetween 0 and 100")
				return
			}

			err = database.CreateWeatherWebhook(m.Author.ID, int64(percent))
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {

				//asd := s.UserChannelCreate(recipientID string)
				//asd, _ := s.UserChannelCreate(m.Author.ID)

				s.ChannelMessageSend(m.ChannelID, "Notification created/updated! You will get notified when the next day has a cloud percentage less than "+fmt.Sprint(percent)+" percent")
			}
		}

		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "drit!")
		}

		if m.Content == "!notify" {
			s.ChannelMessageSend(m.ChannelID, "@"+m.Author.ID+" Hey! Remember to wash your hands")
		}

		if m.Content == "!alljokes" {
			jokes := database.GetAllJokes()
			for _, a := range jokes {
				s.ChannelMessageSend(m.ChannelID, a)
			}
		}

		if m.Content == "!myjokes" {
			jokes := database.GetAllJokesByUserId(m.Author.ID)
			if len(jokes) == 0 {
				s.ChannelMessageSend(m.ChannelID, "You have no jokes yet. Create one with the !createjoke command")
			}
			for _, a := range jokes {
				s.ChannelMessageSend(m.ChannelID, a)
			}
		}

		if strings.HasPrefix(m.Content, "!createjoke ") {
			joke := m.Content[12:]
			err := database.CreateJoke(m.Author.ID, joke)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {
				s.ChannelMessageSend(m.ChannelID, "Joke created")
			}
		}

		// If the message is "pong" reply with "Ping!"
		if m.Content == "pong" {
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		}

		if m.Content == "!userid" {
			s.ChannelMessageSend(m.ChannelID, m.Author.ID)
		}
	*/
}
