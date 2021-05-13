package discordutils

import (
	"fmt"
	"projectGroup23/caching"
	"projectGroup23/discordpkg/constants"

	"github.com/bwmarrin/discordgo"
)

func SendWeatherMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	stringToPrint := constants.GetWeatherStringArray()
	for _, day := range caching.ForecastsCache.Forecasts {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
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

func SendSteamMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	stringToPrint := constants.GetSteamStringArray()
	for _, deal := range caching.DealsCache.Deals {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n%s%s",
			stringToPrint[0], deal.Title, stringToPrint[1], deal.DealID, stringToPrint[2], deal.NormalPrice,
			stringToPrint[3], deal.SalePrice, stringToPrint[4], deal.Savings, stringToPrint[5], deal.MetacriticScore,
			stringToPrint[6], deal.SteamRatingText, stringToPrint[7], deal.SteamRatingPercent, stringToPrint[8], deal.SteamRatingCount))
	}
}

func SendNewsletterMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	stringToPrint := constants.GetNewsletterStringArray()
	for _, article := range caching.NewsCache.Newsletters {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n",
			stringToPrint[0], article.Author, stringToPrint[1], article.Date_published,
			stringToPrint[2], article.Title, stringToPrint[3], article.Description,
			stringToPrint[4], article.Url_to_story))
	}
}

func SendMealplanMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	stringToPrint := constants.GetMealplanMessageArray()
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s\n", stringToPrint[0]))
	for _, meal := range caching.MealsCache.Meals {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s%s\n%s%s\n%s%s\n",
			stringToPrint[1], meal.Title, stringToPrint[2], meal.ReadyInMinutes, stringToPrint[3], meal.Url))
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s\n%s%s\n%s%s\n%s%s\n", stringToPrint[4], fmt.Sprint(caching.MealsCache.Nutrients.Calories),
		stringToPrint[5], fmt.Sprint(caching.MealsCache.Nutrients.Protein), stringToPrint[6], fmt.Sprint(caching.MealsCache.Nutrients.Fat),
		stringToPrint[7], fmt.Sprint(caching.MealsCache.Nutrients.CarboHydrates)))
}
