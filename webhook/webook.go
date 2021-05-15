package webhook

import (
	"fmt"
	"log"
	"projectGroup23/database"
	"projectGroup23/handlers"
	"projectGroup23/utils"
	"time"

	"github.com/bwmarrin/discordgo"
)

// runs at 8am every day
func WebhookRoutine(s *discordgo.Session) {
	webhooks, err := database.GetAllWebhooks()
	if err != nil {
		log.Fatalln(err)
	}
	for _, webhook := range webhooks {
		fmt.Println(webhook)

		/*
			res, err := handlers.WeatherForecastMainHandler(webhook.City)
			if err != nil {
				log.Fatalln(err)
			}
		*/
		fmt.Println("blir denne kjørt?")

		weather, err := handlers.GetWeatherForecastAndIP(webhook.City)
		if err != nil {
			log.Fatalln(err)
		}
		userChannel, _ := s.UserChannelCreate(webhook.UserId)

		stringToPrint := utils.GetWeatherStringArray()
		for _, day := range weather.Forecasts {
			s.ChannelMessageSend(userChannel.ID, utils.WeatherMessageStringFormat(stringToPrint, day))
		}

	}

	timeNow := time.Now()
	var eightAm time.Time
	if timeNow.Hour() < 8 { // if before 8am
		eightAm = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 8, 0, 0, timeNow.Nanosecond(), timeNow.Location()) // 8am today
	} else { // if after 8 am
		eightAm = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day()+1, 8, 0, 0, timeNow.Nanosecond(), timeNow.Location()) // 8am tomorrow
	}

	nextWebhookSeconds := time.Until(eightAm).Seconds()
	time.Sleep(time.Duration(nextWebhookSeconds) * time.Second) // sleeps uintill 8am
	go WebhookRoutine(s)
}
