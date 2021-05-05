package main

import (
	"fmt"
	"os"
	"os/signal"
	"projectGroup23/firebase"
	"projectGroup23/handlers"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {

	firebase.InitFirebase()

	token := "ODM2OTgzNjUyMjUxMzM2Nzc1.YIl7xQ.cuxQXG5lW9Sqmylm6rx4INNiLpc"

	var s, err = discordgo.New("Bot " + token)

	if err != nil {
		panic(err)
	}

	s.AddHandler(messageCreate)
	s.Identify.Intents = discordgo.IntentsGuildMessages

	if err = s.Open(); err != nil {
		panic(err)
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
	// If the message is "ping" reply with "Pong!"
	if m.Content == "!weather" {
		wf := handlers.GetWeatherForecast(1)
		s.ChannelMessageSend(m.ChannelID, "Weather forecast for today")

		for _, day := range wf.Forecasts {
			s.ChannelMessageSend(m.ChannelID, "Date: " + day.Date)
			s.ChannelMessageSend(m.ChannelID, "City: " + day.City)
			s.ChannelMessageSend(m.ChannelID, "Main weather: " + day.Main)
			s.ChannelMessageSend(m.ChannelID, "" + day.Desc)
			s.ChannelMessageSend(m.ChannelID, "Morning temperature: " + fmt.Sprint(day.Morning))
			s.ChannelMessageSend(m.ChannelID, "Day temperature: " + fmt.Sprint(day.Day))
			s.ChannelMessageSend(m.ChannelID, "Evening temperature: " + fmt.Sprint(day.Eve))
			s.ChannelMessageSend(m.ChannelID, "Night temperature: " + fmt.Sprint(day.Night))
			s.ChannelMessageSend(m.ChannelID, "Clouds: " + fmt.Sprint(day.Clouds) + "%")
			s.ChannelMessageSend(m.ChannelID, "Wind: " + fmt.Sprint(day.Wind) + " m/s")
			s.ChannelMessageSend(m.ChannelID, "Probability of precipitation: " + fmt.Sprint(day.POP))
			s.ChannelMessageSend(m.ChannelID, "Rain: " + fmt.Sprint(day.Rain) + " m/s")
			s.ChannelMessageSend(m.ChannelID, "Snow: " + fmt.Sprint(day.Snow) + " m/s")
		}
	}

	if m.Content == "!newsletter" {
		data := handlers.GetDailyNewsLetter()
		for _, article := range data.Newsletters {
			s.ChannelMessageSend(m.ChannelID, "Author: "+article.Author)
			s.ChannelMessageSend(m.ChannelID, "Date: "+article.Date_published)
			s.ChannelMessageSend(m.ChannelID, "Title: "+article.Title)
			s.ChannelMessageSend(m.ChannelID, "Description: "+article.Description)
			s.ChannelMessageSend(m.ChannelID, "Url: "+article.Url_to_story)
			s.ChannelMessageSend(m.ChannelID, " ")

		}
	}

	if m.Content == "!mealplan" {
		mealplan := handlers.GetDailyMealPlan()
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

	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "drit!")
	}

	if m.Content == "!notify" {
		s.ChannelMessageSend(m.ChannelID, "@"+m.Author.ID+" Hey! Remember to wash your hands")
	}

	if m.Content == "!alljokes" {
		jokes := firebase.GetAllJokes()
		for _, a := range jokes {
			s.ChannelMessageSend(m.ChannelID, a)
		}
	}

	if m.Content == "!myjokes" {
		jokes := firebase.GetAllJokesByUserId(m.Author.ID)
		if len(jokes) == 0 {
			s.ChannelMessageSend(m.ChannelID, "You have no jokes yet. Create one with the !createjoke command")
		}
		for _, a := range jokes {
			s.ChannelMessageSend(m.ChannelID, a)
		}
	}

	if strings.HasPrefix(m.Content, "!createjoke ") {
		joke := m.Content[12:]
		err := firebase.CreateJoke(m.Author.ID, joke)
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

}
