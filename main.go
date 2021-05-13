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
	"projectGroup23/handlers"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

var server = "vmdata.database.windows.net"
var port = 1433
var user = "eriksen"
var password = "Tanzania1994!"
var db = "VM_Data"

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, db)

	database.InitFirebase()

	caching.AddCacheModule("cache")

	handlers.GetStoredNewsLetterFromFirestore()
	handlers.GetStoredMealPlannerFromFirestore()
	handlers.GetStoredWeatherForecastFromFirestore()
	handlers.GetStoredSteamDealsFromFirestore()

	token := "ODM2OTgzNjUyMjUxMzM2Nzc1.YIl7xQ.cuxQXG5lW9Sqmylm6rx4INNiLpc"

	var s, err = discordgo.New("Bot " + token)

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println(m.Content)

	switch {
	case m.Content == "!steamdeals":
		dur, _ := time.ParseDuration("10m")
		caching.CacheDeals(m.Content, dur)
		discordutils.SendSteamMessage(s, m)
	case m.Content == "!weather":
		dur, _ := time.ParseDuration("10m")
		caching.CacheForecasts(dur)
		discordutils.SendWeatherMessage(s, m)
	case m.Content == "!mealplan":
		dur, _ := time.ParseDuration("10m")
		caching.CacheMeals(dur)
		discordutils.SendMealplanMessage(s, m)
	case m.Content == "!newsletter":
		dur, _ := time.ParseDuration("10m")
		caching.CacheNews(dur)
		discordutils.SendNewsletterMessage(s, m)
	case m.Content[:5] == "!todo":
		discordutils.SendTodoMessage(s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, "Unable to recognize command, try !help (not implemented) if you need a reminder!")
	}

	/* 	if m.Content == "!steamdeals" {
	   		dur, _ := time.ParseDuration("10m")
	   		caching.CacheDeals(m.Content, dur)
	   		discordutils.SendSteamMessage(s, m)
	   	}

	   	if m.Content == "!weather" {
	   		dur, _ := time.ParseDuration("10m")
	   		caching.CacheForecasts(dur)
	   		discordutils.SendWeatherMessage(s, m)
	   	}

	   	if m.Content == "!mealplan" {
	   		dur, _ := time.ParseDuration("10m")
	   		caching.CacheMeals(dur)
	   		discordutils.SendMealplanMessage(s, m)
	   	}

	   	if m.Content == "!newsletter" {
	   		dur, _ := time.ParseDuration("10m")
	   		caching.CacheNews(dur)
	   		discordutils.SendNewsletterMessage(s, m)
	   	} */

	/* if m.Content == "!todo" {
		allTodos, err := database.GetTodoAll()
		if err != nil {
			log.Fatal("Error reading all todo objects: ", err.Error())
		}
		s.ChannelMessageSend(m.ChannelID, allTodos[1].Description)
	}

	if m.Content[:12] == "!todo create" {
		var todoTask string = m.Content[12:]
		fmt.Println(todoTask)
		if len(todoTask) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Missing description for todo task!")
			return
		}

		err := database.CreateTodoObject(todoObject)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong while creating todo object")
			fmt.Println(err)
		}
	}

	if m.Content == "!todo mylist" {
		allTodos, err := database.GetTodoObject(m.Author.ID)
		if err != nil {
			log.Fatal("Error reading all todo objects: ", err.Error())
		}
		fmt.Println(allTodos)
		//s.ChannelMessageSend(m.ChannelID, allTodos[1].Description)

		for i, todo := range allTodos {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint(i+1)+": "+todo.Description+" -- status: "+todo.State)
		}
	}

	if m.Content[:12] == "!todo delete" {
		var todoId string = m.Content[12:]
		fmt.Println(todoId)
		if len(todoId) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Missing id for todo task!")
			return
		}
		conv, err := strconv.Atoi(todoId[1:])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Id needs to be a number!")
			return
		}

		res, err := convertIndexToId((conv - 1), m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		err = database.DeleteTodoObject(res)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
	}

	if m.Content[:12] == "!todo update" {
		var argString string = m.Content[12:]
		fmt.Println(argString)
		if len(argString) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Missing id for todo task!")
			return
		}

		args := strings.Fields(argString)

		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Wrong format!")
			return
		}

		fmt.Println(args)

		conv, err := strconv.Atoi(args[0])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Id needs to be a number!")
			return
		}

		res, err := convertIndexToId((conv - 1), m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		fmt.Println(res)

		var updateTodo string

		for i, word := range args[1:] {
			updateTodo += word
			if i != len(args[1:])-1 {
				updateTodo += " "
			}
		}

		err = database.UpdateTodoObject(res, updateTodo)
		if err != nil {
			fmt.Println(err)
		}

	}

	// If the message is "ping" reply with "Pong!"


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
