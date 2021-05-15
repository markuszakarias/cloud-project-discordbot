package discord

import (
	"errors"
	"fmt"
	"log"
	"os"
	"projectGroup23/caching"
	"projectGroup23/database"
	"projectGroup23/handlers"
	"projectGroup23/structs"
	"projectGroup23/utils"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func DiagMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Get api key from env variable
	apikey := os.Getenv("MEALS_KEY")
	url := "https://api.spoonacular.com/mealplanner/generate?timeFrame=day&apiKey=" + apikey
	statusCodeMeal := utils.CheckStatusCodeApi(url)

	country := "norway"
	apikey = os.Getenv("NEWS_KEY")
	url = "https://newsapi.org/v2/top-headlines?country=" + country + "&apiKey=" + apikey
	statusCodeNews := utils.CheckStatusCodeApi(url)

	url = "https://www.cheapshark.com/api/1.0/deals"
	statusCodeSteam := utils.CheckStatusCodeApi(url)

	location := "lillehammer"
	apikey = os.Getenv("WEATHER_KEY")
	url = "https://api.openweathermap.org/data/2.5/forecast/daily?q=" + location + "&units=metric&cnt=1&appid=" + apikey
	statusCodeWeather := utils.CheckStatusCodeApi(url)

	url = "https://restcountries.eu/rest/v2/name/" + country + "?fullText=true"
	statusCodeRestCountries := utils.CheckStatusCodeApi(url)

	url = "https://ipwhois.app/json/"
	statusCodeIp := utils.CheckStatusCodeApi(url)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
		"Meals API status code: %s\nNews API status code: %s\nSteam deals API status code: %s\nWeather API status code: %s\nIP location API status code: %s\nRest countries API status code: %s", statusCodeMeal, statusCodeNews, statusCodeSteam, statusCodeWeather, statusCodeRestCountries, statusCodeIp))
}

// SendWeatherMessage - sends appropriate message to the bot based on command and parameters
func SendWeatherMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {

	// Sets the default location to the GeoPosition of the computer making the request
	defaultLocation, err := utils.GetIPLocation()
	if err != nil {
		return err
	}

	// Sets duration of cache
	dur, _ := time.ParseDuration("20s")

	// Splits up the input command(s)
	str := strings.Fields(m.Content)
	fmt.Println(str)

	// Checks if there are any parameters
	if len(str) < 2 {
		err := caching.CacheForecasts(defaultLocation, dur)
		if err != nil {
			return err
		}
	} else {
		location := strings.Title(strings.ToLower(str[1]))
		err := caching.CacheForecasts(location, dur)
		if err != nil {
			return err
		}
	}

	// Printer function
	stringToPrint := utils.GetWeatherStringArray()
	for _, day := range caching.ForecastsCache.Forecasts {
		s.ChannelMessageSend(m.ChannelID, utils.WeatherMessageStringFormat(stringToPrint, day))
	}
	return nil
}

// SendSteamMessage - sends appropriate message to the bot based on command and parameters
func SendSteamMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {

	// count variable to set amount to be displayed in the output
	var count int
	var err error

	// Splits up the command(s)
	str := strings.Fields(m.Content)

	dur, _ := time.ParseDuration("20s")

	// Checks if there are any parameters
	if len(str) < 2 {
		count = 5
		err = caching.CacheDeals(m.Content, dur)
		if err != nil {
			return err
		}
	} else {
		// If there are any number parameters, convert
		count, err = strconv.Atoi(str[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Id needs to be a number!")
			return err
		}
		if count < 1 || count > 4 {
			return errors.New("number of steamdeals needs to be between 1-4")
		}
		err = caching.CacheDeals(m.Content, dur)
		if err != nil {
			return err
		}
	}

	// Printer function
	stringToPrint := utils.GetSteamStringArray()
	for i := 0; i < count; i++ {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n",
			stringToPrint[0], caching.DealsCache.Deals[i].Title,
			stringToPrint[1], caching.DealsCache.Deals[i].DealID,
			stringToPrint[2], caching.DealsCache.Deals[i].NormalPrice,
			stringToPrint[3], caching.DealsCache.Deals[i].SalePrice,
			stringToPrint[4], caching.DealsCache.Deals[i].Savings))
	}
	return nil
}

// SendNewsletterMessage - sends appropriate message to the bot based on command and parameters
func SendNewsletterMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {

	// Default amount of newsletter articles
	var count = 3

	// Default country to retrieve newsletter
	var country = "no"

	// Sets duration of cache
	dur, _ := time.ParseDuration("20s")

	// split up the parameter for validation and passing data
	str := strings.Fields(m.Content)

	// Checks if there are any parameters
	if len(str) < 2 {
		err := caching.CacheNews(country, dur)
		if err != nil {
			return err
		}
	} else if len(str) > 1 {
		alpha2Code, err := utils.Get2AlphaCode(str[1])
		if err != nil {
			return err
		}

		if len(str) > 2 {
			count, err = strconv.Atoi(str[2])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Id needs to be a number!")
				return err
			}
			if count < 1 || count > 3 {
				return errors.New("number of articles needs to be between 1-3")
			}
		}

		result := strings.ToLower(alpha2Code)
		err = caching.CacheNews(result, dur)
		if err != nil {
			return err
		}
	}

	// Print function
	stringToPrint := utils.GetNewsletterStringArray()
	if len(caching.NewsCache.Newsletters) < count {
		count = len(caching.NewsCache.Newsletters)
	}

	for i := 0; i < count; i++ {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n",
			stringToPrint[0], caching.NewsCache.Newsletters[i].Author,
			stringToPrint[1], caching.NewsCache.Newsletters[i].Date_published,
			stringToPrint[2], caching.NewsCache.Newsletters[i].Title,
			stringToPrint[3], caching.NewsCache.Newsletters[i].Description,
			stringToPrint[4], caching.NewsCache.Newsletters[i].UrlToStory))
	}
	return nil
}

// SendMealplanMessage - sends appropriate message to the bot based on command and parameters
func SendMealplanMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	dur, _ := time.ParseDuration("20s")
	err := caching.CacheMeals(dur)
	if err != nil {
		return err
	}
	// Printer function
	stringToPrint := utils.GetMealplanMessageArray()
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s\n", stringToPrint[0]))
	for _, meal := range caching.MealsCache.Meals {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s%s\n%s%s\n%s%s\n",
			stringToPrint[1], meal.Title, stringToPrint[2], meal.ReadyInMinutes, stringToPrint[3], meal.Url))
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s\n%s%s\n%s%s\n%s%s\n", stringToPrint[4], fmt.Sprint(caching.MealsCache.Nutrients.Calories),
		stringToPrint[5], fmt.Sprint(caching.MealsCache.Nutrients.Protein), stringToPrint[6], fmt.Sprint(caching.MealsCache.Nutrients.Fat),
		stringToPrint[7], fmt.Sprint(caching.MealsCache.Nutrients.CarboHydrates)))
	return nil
}

// SendTodoMessage - sends appropriate message to the bot based on command and parameters
func SendTodoMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {

	// strings to be assigned value
	var createTodo string
	var updateTodo string

	// struct to be used for todo object
	var todoObject structs.TodoStruct
	todoObject.Userid = m.Author.ID
	todoObject.State = "active"

	// split up command(s)
	str := strings.Fields(m.Content)
	fmt.Println(str)

	// Checks if there are any parameters with todo command
	if len(str) < 2 {
		return errors.New("command missing for !todo ")
	}

	// Switch to handle the different parameters
	switch {
	case str[1] == "mylist":
		// gets the todo objects from azure sql
		allTodos, err := database.GetTodoObject(m.Author.ID)
		if err != nil {
			log.Fatal("Error reading all todo objects: ", err.Error())
		}
		for i, todo := range allTodos {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint(i+1)+": "+todo.Description+" -- status: "+todo.State)
		}
	case str[1] == "create":
		// splits up the parameters, taking everything after the "create" parameter as the todo task
		createTodo = strings.Join(str[2:], " ")

		if createTodo == "" {
			return errors.New("missing description for todo task")
		}
		// inserts it into the struct object
		todoObject.Description = createTodo
		// inserts it into azure sql
		err := database.CreateTodoObject(todoObject)
		if err != nil {
			return errors.New("something went wrong while creating todo object")
		}
		s.ChannelMessageSend(m.ChannelID, "Task was created.")
	case str[1] == "delete":
		if str[2] == "" {
			return errors.New("missing id for todo task")
		}
		// converts string number from parameter to int number
		// will exit if failure
		conv, err := strconv.Atoi(str[2])
		if err != nil {
			return errors.New("id needs to be a number")
		}

		// converts the id number to the id on azure sql
		res, err := database.ConvertIndexToId((conv - 1), m.Author.ID)
		if err != nil {
			return err
		}
		err = database.DeleteTodoObject(res)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, "Task with id: "+fmt.Sprint(conv)+" was deleted.")
	case str[1] == "update":
		if str[2] == "" {
			return errors.New("missing id for todo task")
		}
		if str[3] == "" {
			return errors.New("missing data to update")
		}

		// converts string number from parameter to int number
		// will exit if failure
		conv, err := strconv.Atoi(str[2])
		if err != nil {
			return errors.New("id needs to be a number")
		}

		// converts the id number to the id on azure sql
		res, err := database.ConvertIndexToId((conv - 1), m.Author.ID)
		if err != nil {
			return err
		}

		// splits up the parameters and takes everything from index 3 as a string
		updateTodo = strings.Join(str[3:], " ")

		// updates the object in azure sql
		err = database.UpdateTodoObject(res, updateTodo)
		if err != nil {
			return err
		}

		s.ChannelMessageSend(m.ChannelID, "Task with id: "+fmt.Sprint(conv)+" was updated.")
	case str[1] == "finished" || str[1] == "inactive" || str[1] == "active":
		if str[2] == "" {
			return errors.New("missing id for todo task")
		}

		// converts string number from parameter to int number
		// will exit if failure
		conv, err := strconv.Atoi(str[2])
		if err != nil {
			return errors.New("id needs to be a number")
		}

		// converts the id number to the id on azure sql
		res, err := database.ConvertIndexToId((conv - 1), m.Author.ID)
		if err != nil {
			return err
		}

		// updates the object in azure sql
		err = database.UpdateTodoObjectStatus(res, str[1])
		if err != nil {
			return err
		}

		s.ChannelMessageSend(m.ChannelID, "Task with id: "+fmt.Sprint(conv)+" status was updated with: "+fmt.Sprint(str[1]))
	}
	return nil
}

// SendHelpMessage - sends appropriate message to the bot based on command and parameters
func SendHelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// split up command(s)
	str := strings.Fields(m.Content)

	if len(str) < 2 {
		stringToPrint := utils.GetHelpMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n%s\n%s\n%s\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n",
			stringToPrint[0], stringToPrint[1], stringToPrint[2],
			stringToPrint[3], stringToPrint[4], stringToPrint[5],
			stringToPrint[6], stringToPrint[7], stringToPrint[8],
			stringToPrint[9]))
		return
	}

	switch {
	case str[1] == "todo":
		stringToPrint := utils.GetHelpTodoMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n%s\n%s",
			stringToPrint[0], stringToPrint[1], stringToPrint[2],
			stringToPrint[3], stringToPrint[4], stringToPrint[5],
			stringToPrint[6], stringToPrint[7], stringToPrint[8],
			stringToPrint[9], stringToPrint[10], stringToPrint[11],
			stringToPrint[12], stringToPrint[13], stringToPrint[14],
			stringToPrint[15]))
		return
	case str[1] == "weather":
		stringToPrint := utils.GetHelpWeatherMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n%s%s%s\n\n%s\n%s\n",
			stringToPrint[0], stringToPrint[1], stringToPrint[3],
			stringToPrint[1], stringToPrint[4], stringToPrint[5]))
		return
	case str[1] == "newsletter":
		stringToPrint := utils.GetHelpNewsletterMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n%s%s%s\n\n%s\n%s\n",
			stringToPrint[0], stringToPrint[1], stringToPrint[3],
			stringToPrint[1], stringToPrint[4], stringToPrint[5]))
		return
	case str[1] == "steamdeals":
		stringToPrint := utils.GetHelpSteamdealsMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n%s%s%s\n\n%s\n%s\n",
			stringToPrint[0], stringToPrint[1], stringToPrint[3],
			stringToPrint[1], stringToPrint[4], stringToPrint[5]))
		return
	case str[1] == "mealplan":
		stringToPrint := utils.GetHelpMealplanMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n%s%s%s\n\n%s\n",
			stringToPrint[0], stringToPrint[1], stringToPrint[3],
			stringToPrint[1], stringToPrint[4]))
		return
	case str[1] == "notifyweather":
		stringToPrint := utils.GetHelpNotifyWeatherMessageArray()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"%s\n\n%s%s%s\n\n%s\n\n%s\n%s\n%s\n",
			stringToPrint[0], stringToPrint[1], stringToPrint[3],
			stringToPrint[1], stringToPrint[4], stringToPrint[5],
			stringToPrint[6], stringToPrint[7]))
		return
	}
}

// NotifyWeather - sends appropriate message to the bot based on command and parameters
func NotifyWeather(s *discordgo.Session, m *discordgo.MessageCreate) error {

	// split up the command(s)
	str := strings.Fields(m.Content)
	fmt.Println(str)

	// Checks if there are any parameters with the command
	if len(str) < 2 {
		return errors.New("missing city name")
	}

	// Removes webhook
	if str[1] == "remove" {
		database.DeleteWeatherWebhook(m.Author.ID)
		s.ChannelMessageSend(m.ChannelID, "Notification is removed")
		return nil
	}

	// Creates a weather API call
	_, err := handlers.WeatherForecastMainHandler(str[1])
	if err != nil {
		return err
	}

	// create a webhook on the weather
	err = database.CreateWeatherWebhook(m.Author.ID, str[1])

	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, "Notification is registered. You will be notified with the weather information at 8 am")

	return nil
}

// SendJokeMessage - sends appropriate message to the bot based on command and parameters
func SendJokeMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {

	// split up command(s)
	str := strings.Fields(m.Content)
	switch {
	case len(str) == 1: // !joke
		joke, err := database.GetRandomJoke()
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, joke)
		return nil
	case len(str) == 2 && str[1] == "myjokes":
		jokes, err := database.GetAllJokesByUserId(m.Author.ID)
		if err != nil {
			return err
		}
		allJokeString := ""
		for i, joke := range jokes {
			allJokeString += fmt.Sprint(i) + ". " + joke + "\n"
		}
		if allJokeString == "" {
			s.ChannelMessageSend(m.ChannelID, "You have not created any jokes yet")
		} else {
			s.ChannelMessageSend(m.ChannelID, allJokeString)
		}
		return nil
	case len(str) == 2 && str[1] == "create": // misses joke text
		return errors.New("missing joke text")
	case len(str) > 2 && str[1] == "create": // !joke create text here
		joke := strings.Join(str[2:], " ")
		database.CreateJoke(m.Author.ID, joke)
		s.ChannelMessageSend(m.ChannelID, "joke created")
		return nil
	default:
		return errors.New("something is wrong with your joke command")
	}

}
