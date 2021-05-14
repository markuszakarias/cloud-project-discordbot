package webhook

import (
	"fmt"
	"log"
	"projectGroup23/database"
	"projectGroup23/discordpkg/constants"
	"projectGroup23/handlers"
	"time"

	"github.com/bwmarrin/discordgo"
)

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

		stringToPrint := constants.GetWeatherStringArray()
		for _, day := range weather.Forecasts {
			s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
				"%s%s\n %s%s\n %s%s%s\n %s%s%s\n %s%s%s\n %s%s\n %s%s%s\n %s%s%s\n %s\n %s%s%s\n %s%s%s\n %s%s%s\n %s%s%s\n",
				stringToPrint[0], day.Date, stringToPrint[1], day.City, day.Main,
				stringToPrint[2], day.Desc, stringToPrint[3], fmt.Sprint(day.Clouds),
				stringToPrint[15], stringToPrint[4], fmt.Sprint(day.Wind), stringToPrint[5],
				stringToPrint[7], fmt.Sprint(day.POP), stringToPrint[8], fmt.Sprint(day.Rain), stringToPrint[5],
				stringToPrint[9], fmt.Sprint(day.Snow), stringToPrint[5], stringToPrint[10],
				stringToPrint[11], fmt.Sprint(day.Morning), stringToPrint[6],
				stringToPrint[12], fmt.Sprint(day.Day), stringToPrint[6],
				stringToPrint[13], fmt.Sprint(day.Eve), stringToPrint[6],
				stringToPrint[14], fmt.Sprint(day.Night), stringToPrint[6]))
		}

	}
	time.Sleep(time.Duration(9999) * time.Second) // waits 15 minutes
	go WebhookRoutine(s)
}