package utils

// GetWeatherStringArray returns an immutable array of the messages that make up !weather
func GetWeatherStringArray() []string {
	return []string{":calendar: ", ":map: ", " - ", ":cloud: ", ":dash: ",
		" m/s", " Celsius", "Probability of precipitation: ", ":cloud_rain: ", ":cloud_snow: ",
		"Temperature:", ":city_sunrise: ", ":cityscape: ", ":city_dusk: ", ":night_with_stars: ", "%"}
}

// GetSteamStringArray returns an immutable array of the messages that make up !steamdeals
func GetSteamStringArray() []string {
	return []string{":video_game: ",
		"**Deal ID:** ",
		"**Normal Price:** ",
		"**Sale Price:** ",
		"**Savings:** "}
}

// GetNewsletterStringArray returns an immutable arrray of the messages that make up !newsletter
func GetNewsletterStringArray() []string {
	return []string{"Author: ", "Date: ", "Title: ", "Description: ", "Url: "}
}

// GetMealplanMessageArray returns an arrray of the messages that make up !mealplan
func GetMealplanMessageArray() []string {
	return []string{"Meal message: ", "Title: ", "Ready in minutes: ", "url: ", "Calories: ", "Protein: ", "Fat: ", "CarboHydrates: "}
}

// GetHelpMessageArray returns an arrray of the messages that format the base helper message
func GetHelpMessageArray() []string {
	return []string{"Helper message for discord bot",
		"With this discord bot you can retrieve information about news, weather",
		"and mealplans from different APIs. You can also create todo tasks that are",
		"attached to you discord ID. To get more information about the different",
		"commands, use the help command witht the appropriate parameter:",
		"**Todo**              --> !help todo",
		"**Newsletter**   --> !help newsletter",
		"**Weather**       --> !help weather",
		"**Mealplan**      --> !help mealplan",
		"**Steamdeals**  --> !help steamdeals",
	}
}

// GetHelpTodoMessageArray returns an arrray of the messages that format help message for todo
func GetHelpTodoMessageArray() []string {
	return []string{"Helper message for discord bot",
		"Todo list :pencil:",
		"With the todo list you can view, create, update and delete todo tasks.",
		"They are attached to your discord ID. Examples:",
		"**View your todo list**",
		"!todo mylist",
		"**Create todo task**",
		"!todo create this is my task",
		"**Update todo task**",
		"!todo update <taskid> this is the updated task",
		"**Label todo task as finished**",
		"!todo update <taskid> finished",
		"**Label todo task as inactive**",
		"todo update <taskid> inactive",
		"**Delete todo task**",
		"!todo delete <taskid>",
	}
}

// GetHelpWeatherMessageArray returns an arrray of the messages that make up !help weather
func GetHelpWeatherMessageArray() []string {
	return []string{"Helper message for Weather :partly_sunny:",
		"`",
		"**",
		"!weather",
		"Using the weather command you can call with just !weather to get your location!",
		"If you supply a parameter like: !weather Fredrikstad, you will get the weather for that location"}
}

// GetHelpNotifyWeatherMessageArray returns an arrray of the messages that make up !help notifyweather
func GetHelpNotifyWeatherMessageArray() []string {
	return []string{"Helper message for Notify weather :partly_sunny: :grey_exclamation:",
		"`",
		"**",
		"!notifyweather",
		"Using the weather notify command to register a webhook!",
		"This will subscribe you to get a message every day at 08.00 am with the current days weather report.",
		"You can only be subscribed to one city at a time, so calling a new notify weather command with a different",
		"city will replace the old one."}
}

// GetHelpNewsletterMessageArray returns an arrray of the messages that make up !help newsletter
func GetHelpNewsletterMessageArray() []string {
	return []string{"Helper message for Newsletter :earth_africa:",
		"`",
		"**",
		"!newsletter",
		"call the command to get current 3 norwegian headlines.",
		"You can also supply a country and headlines between 1-3 such as: `!newsletter Sweden 2`"}
}

// GetHelpMealplanMessageArray returns an arrray of the messages that make up !help mealplan
func GetHelpMealplanMessageArray() []string {
	return []string{"Helper message for Mealplan :ramen:",
		"`",
		"**",
		"!mealplan",
		"run the command to get a breakfeast -> dinner -> snack mealplan for today!"}
}

// GetHelpSteamdealsMessageArray returns an arrray of the messages that make up !help steamdeals
func GetHelpSteamdealsMessageArray() []string {
	return []string{"Helper message for Steamdeals :video_game:",
		"`",
		"**",
		"!steamdeals",
		"call the command to get a list of 4 current deals on the steam store",
		"if you can less you can specify an amount: `!steamdeals 2`"}
}
