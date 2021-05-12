package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"projectGroup23/caching"
	"projectGroup23/database"
	"syscall"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"

	"github.com/bwmarrin/discordgo"
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

	caching.GetCachedNewsLetterFromFirestore()
	caching.GetCachedMealPlannerFromFirestore()
	caching.GetCachedWeatherForecastFromFirestore()
	caching.GetSteamdealsFromFirestore()

	token := "ODM2OTgzNjUyMjUxMzM2Nzc1.YIl7xQ.cuxQXG5lW9Sqmylm6rx4INNiLpc"

	var s, err = discordgo.New("Bot " + token)

	if err = s.Open(); err != nil {
		panic(err)
	}

	go database.WebhookRoutine(s)
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

func convertIndexToId(i int, userid string) (int, error) {
	resp, err := database.GetTodoObject(userid)
	if err != nil {
		return 0, err
	}
	if len(resp) <= i {
		return 0, errors.New("id does not exist")
	}
	deleteId := resp[i].Id

	return deleteId, nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	/* if m.Content == "!mealplan" {
		mealplan := caching.MealPlannerTest()
		s.ChannelMessageSend(m.ChannelID, "meal message: "+mealplan.MealMessage)
		for _, meal := range mealplan.Meals {
			s.ChannelMessageSend(m.ChannelID, "title: "+meal.Title)
			s.ChannelMessageSend(m.ChannelID, "ready in minuts: "+meal.ReadyInMinutes)
			s.ChannelMessageSend(m.ChannelID, "url: "+meal.Url)
		}
		s.ChannelMessageSend(m.ChannelID, "Calories: "+fmt.Sprint(mealplan.Nutrients.Calories))
		s.ChannelMessageSend(m.ChannelID, "Protein: "+fmt.Sprint(mealplan.Nutrients.Protein))
		s.ChannelMessageSend(m.ChannelID, "Fat: "+fmt.Sprint(mealplan.Nutrients.Fat))
		s.ChannelMessageSend(m.ChannelID, "CarboHydrates: "+fmt.Sprint(mealplan.Nutrients.CarboHydrates))

	} */

	/* if m.Content == "!steamdeals" {
		deals := caching.SteamdealsTest(m.Content)
		s.ChannelMessageSend(m.ChannelID, "Here are your steam deal(s): ")
		for _, deal := range deals.Deals {
			s.ChannelMessageSend(m.ChannelID, "title: "+deal.Title)
			s.ChannelMessageSend(m.ChannelID, "DealID: "+deal.DealID)
			s.ChannelMessageSend(m.ChannelID, "NormalPrice: "+deal.NormalPrice)
			s.ChannelMessageSend(m.ChannelID, "SalePrice: "+deal.SalePrice)
			s.ChannelMessageSend(m.ChannelID, "Savings: "+deal.Savings)
			s.ChannelMessageSend(m.ChannelID, "MetacriticScore: "+deal.MetacriticScore)
			s.ChannelMessageSend(m.ChannelID, "SteamRatingText: "+deal.SteamRatingText)
			s.ChannelMessageSend(m.ChannelID, "SteamRatingPercent: "+deal.SteamRatingPercent)
			s.ChannelMessageSend(m.ChannelID, "SteamRatingCount: "+deal.SteamRatingCount)
		}
	} */

	/* if m.Content == "!newsletter" {
		caching.AddCacheModule("news")
		dur, _ := time.ParseDuration("10m")
		caching.CacheNews(dur)

		for _, article := range caching.NewsCache.Newsletters {
			s.ChannelMessageSend(m.ChannelID, "Author: "+article.Author)
			s.ChannelMessageSend(m.ChannelID, "Date: "+article.Date_published)
			s.ChannelMessageSend(m.ChannelID, "Title: "+article.Title)
			s.ChannelMessageSend(m.ChannelID, "Description: "+article.Description)
			s.ChannelMessageSend(m.ChannelID, "Url: "+article.Url_to_story)
			s.ChannelMessageSend(m.ChannelID, " ")
		}
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

		todoTask = todoTask[1:]

		var todoObject structs.Todo_struct
		todoObject.Userid = m.Author.ID
		todoObject.Description = todoTask
		todoObject.State = "active"

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

	if m.Content == "!steamdeals" {
		deals := handlers.GetSteamDeals(m.Content)
		s.ChannelMessageSend(m.ChannelID, "Here are your steam deal(s): ")
		for _, deal := range deals.Deals {
			s.ChannelMessageSend(m.ChannelID, "title: "+deal.Title)
			s.ChannelMessageSend(m.ChannelID, "DealID: "+deal.DealID)
			s.ChannelMessageSend(m.ChannelID, "NormalPrice: "+deal.NormalPrice)
			s.ChannelMessageSend(m.ChannelID, "SalePrice: "+deal.SalePrice)
			s.ChannelMessageSend(m.ChannelID, "Savings: "+deal.Savings)
			s.ChannelMessageSend(m.ChannelID, "MetacriticScore: "+deal.MetacriticScore)
			s.ChannelMessageSend(m.ChannelID, "SteamRatingText: "+deal.SteamRatingText)
			s.ChannelMessageSend(m.ChannelID, "SteamRatingPercent: "+deal.SteamRatingPercent)
			s.ChannelMessageSend(m.ChannelID, "SteamRatingCount: "+deal.SteamRatingCount)
		}
	} */
}
