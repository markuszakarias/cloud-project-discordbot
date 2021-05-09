package main

import (
	"fmt"
	"os"
	"os/signal"
	"projectGroup23/firebase"
	"projectGroup23/handlers"
	"strconv"
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

	go firebase.WebhookRoutine(s)
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

		for _, day := range wf.Forecasts {
			s.ChannelMessageSend(m.ChannelID, ":calendar: "+day.Date)
			s.ChannelMessageSend(m.ChannelID, ":map: "+day.City)
			s.ChannelMessageSend(m.ChannelID, day.Main+" - "+day.Desc)
			s.ChannelMessageSend(m.ChannelID, ":cloud: "+fmt.Sprint(day.Clouds)+"%")
			s.ChannelMessageSend(m.ChannelID, ":dash: "+fmt.Sprint(day.Wind)+" m/s")
			s.ChannelMessageSend(m.ChannelID, "Probability of precipitation: "+fmt.Sprint(day.POP))
			s.ChannelMessageSend(m.ChannelID, ":cloud_rain: "+fmt.Sprint(day.Rain)+" m/s")
			s.ChannelMessageSend(m.ChannelID, ":cloud_snow: "+fmt.Sprint(day.Snow)+" m/s")
			s.ChannelMessageSend(m.ChannelID, "Temperature:")
			s.ChannelMessageSend(m.ChannelID, ":city_sunrise: "+fmt.Sprint(day.Morning)+" Celsius")
			s.ChannelMessageSend(m.ChannelID, ":cityscape: "+fmt.Sprint(day.Day)+" Celsius")
			s.ChannelMessageSend(m.ChannelID, ":city_dusk: "+fmt.Sprint(day.Eve)+" Celsius")
			s.ChannelMessageSend(m.ChannelID, ":night_with_stars: "+fmt.Sprint(day.Night)+" Celsius")
		}
	}

	if m.Content == "!notifyweather remove" {
		err := firebase.DeleteWebhook(m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "notification removed!")
			return
		}
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

		err = firebase.CreateWeatherWebhook(m.Author.ID, int64(percent))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		} else {

			//asd := s.UserChannelCreate(recipientID string)
			//asd, _ := s.UserChannelCreate(m.Author.ID)

			s.ChannelMessageSend(m.ChannelID, "Notification created/updated! You will get notified when the next day has a cloud percentage less than "+fmt.Sprint(percent)+" percent")
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
