package discordutils

import (
	"errors"
	"fmt"
	"log"
	"projectGroup23/caching"
	"projectGroup23/database"
	"projectGroup23/discordpkg/constants"
	"projectGroup23/handlers"
	"projectGroup23/structs"
	"projectGroup23/utils"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SendWeatherMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {

	defaultLocation, err := utils.GetIPLocation()
	if err != nil {
		return err
	}

	str := strings.Fields(m.Content)
	fmt.Println(str)

	if len(str) < 2 {
		dur, _ := time.ParseDuration("20s")
		err := caching.CacheForecasts(defaultLocation, dur)
		if err != nil {
			return err
		}
	} else {
		location := strings.Title(strings.ToLower(str[1]))
		dur, _ := time.ParseDuration("10m")
		err := caching.CacheForecasts(location, dur)
		if err != nil {
			return err
		}
	}

	stringToPrint := constants.GetWeatherStringArray()
	for _, day := range caching.ForecastsCache.Forecasts {
		s.ChannelMessageSend(m.ChannelID, utils.WeatherMessageStringFormat(stringToPrint, day))
	}
	return nil
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

func SendTodoMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	var createTodo string
	var updateTodo string

	var todoObject structs.Todo_struct
	todoObject.Userid = m.Author.ID
	todoObject.State = "active"

	str := strings.Fields(m.Content)
	fmt.Println(str)

	if len(str) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Command missing for !todo. ")
		return
	}

	switch {
	case str[1] == "mylist":
		allTodos, err := database.GetTodoObject(m.Author.ID)
		if err != nil {
			log.Fatal("Error reading all todo objects: ", err.Error())
		}
		for i, todo := range allTodos {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint(i+1)+": "+todo.Description+" -- status: "+todo.State)
		}
	case str[1] == "create":
		createTodo = strings.Join(str[2:], " ")

		if createTodo == "" {
			s.ChannelMessageSend(m.ChannelID, "Missing description for todo task!")
			return
		}

		todoObject.Description = createTodo

		err := database.CreateTodoObject(todoObject)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong while creating todo object")
			fmt.Println(err)
		}
		s.ChannelMessageSend(m.ChannelID, "Task was created.")
	case str[1] == "delete":
		if str[2] == "" {
			s.ChannelMessageSend(m.ChannelID, "Missing id for todo task!")
			return
		}
		conv, err := strconv.Atoi(str[2])
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
		s.ChannelMessageSend(m.ChannelID, "Task with id: "+fmt.Sprint(conv)+" was deleted.")
	case str[1] == "update":
		if str[2] == "" {
			s.ChannelMessageSend(m.ChannelID, "Missing id for todo task!")
			return
		}
		if str[3] == "" {
			s.ChannelMessageSend(m.ChannelID, "Missing data to update!")
			return
		}

		conv, err := strconv.Atoi(str[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Id needs to be a number!")
			return
		}

		res, err := convertIndexToId((conv - 1), m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		updateTodo = strings.Join(str[3:], " ")

		err = database.UpdateTodoObject(res, updateTodo)
		if err != nil {
			fmt.Println(err)
		}

		s.ChannelMessageSend(m.ChannelID, "Task with id: "+fmt.Sprint(conv)+" was updated.")
	case str[1] == "finished" || str[1] == "inactive" || str[1] == "active":
		if str[2] == "" {
			s.ChannelMessageSend(m.ChannelID, "Missing id for todo task!")
			return
		}
		conv, err := strconv.Atoi(str[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Id needs to be a number!")
			return
		}

		res, err := convertIndexToId((conv - 1), m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		err = database.UpdateTodoObjectStatus(res, str[1])
		if err != nil {
			fmt.Println(err)
		}

		s.ChannelMessageSend(m.ChannelID, "Task with id: "+fmt.Sprint(conv)+" status was updated with: "+fmt.Sprint(str[1]))
	}
}

func SendHelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	str := strings.Fields(m.Content)
	fmt.Println(str)

	if len(str) < 2 {
		stringToPrint := constants.GetHelpMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n%s\n%s\n%s\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n",
			stringToPrint[0], stringToPrint[1], stringToPrint[2],
			stringToPrint[3], stringToPrint[4], stringToPrint[5],
			stringToPrint[6], stringToPrint[7], stringToPrint[8],
			stringToPrint[9]))
		return
	}

	if str[1] == "todo" {
		stringToPrint := constants.GetHelpTodoMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n%s\n%s",
			stringToPrint[0], stringToPrint[1], stringToPrint[2],
			stringToPrint[3], stringToPrint[4], stringToPrint[5],
			stringToPrint[6], stringToPrint[7], stringToPrint[8],
			stringToPrint[9], stringToPrint[10], stringToPrint[11],
			stringToPrint[12], stringToPrint[13], stringToPrint[14],
			stringToPrint[15]))
		return
	}
	if str[1] == "weather" {
		stringToPrint := constants.GetHelpWeatherMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n",
			stringToPrint[0]))
		return

	}
	if str[1] == "newsletter" {
		stringToPrint := constants.GetHelpNewsletterMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n",
			stringToPrint[0]))
		return

	}
	if str[1] == "steamdeals" {
		stringToPrint := constants.GetHelpSteamdealsMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n",
			stringToPrint[0]))
		return

	}
	if str[1] == "mealplan" {
		stringToPrint := constants.GetHelpMealplanMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n",
			stringToPrint[0]))
		return

	}

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

func NotifyWeather(s *discordgo.Session, m *discordgo.MessageCreate) error {

	str := strings.Fields(m.Content)
	fmt.Println(str)

	if len(str) < 2 {
		return errors.New("Missing city name")
	}

	_, err := handlers.WeatherForecastMainHandler(str[1])
	if err != nil {
		return err
	}

	err = database.CreateWeatherWebhook(m.Author.ID, str[1])

	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, "Notification is registred. You will be notified with the weather information at 8 am")

	return nil
}
