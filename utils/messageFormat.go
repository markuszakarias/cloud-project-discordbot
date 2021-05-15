package utils

func GetWeatherStringArray() []string {
	return []string{":calendar: ", ":map: ", " - ", ":cloud: ", ":dash: ",
		" m/s", " Celsius", "Probability of precipitation: ", ":cloud_rain: ", ":cloud_snow: ",
		"Temperature:", ":city_sunrise: ", ":cityscape: ", ":city_dusk: ", ":night_with_stars: ", "%"}
}

func GetSteamStringArray() []string {
	return []string{":video_game: ",
		"**Deal ID:** ",
		"**Normal Price:** ",
		"**Sale Price:** ",
		"**Savings:** "}
}

func GetNewsletterStringArray() []string {
	return []string{"Author: ", "Date: ", "Title: ", "Description: ", "Url: "}
}

func GetMealplanMessageArray() []string {
	return []string{"Meal message: ", "Title: ", "Ready in minutes: ", "url: ", "Calories: ", "Protein: ", "Fat: ", "CarboHydrates: "}
}

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

func GetHelpWeatherMessageArray() []string {
	return []string{"Helper message for Weather :partly_sunny:",
		"`",
		"**",
		"Using the weather command you can call with just !weather to get your location!",
		"If you supply a paramter like: !weather Fredrikstad, you will get the weather for that location"}
}

func GetHelpNewsletterMessageArray() []string {
	return []string{"Helper message for Newsletter :earth_africa:",
		"`",
		"**",
		"call !newsletter to get current 3 norwegian headlines",
		"You can also supply a country and headlines between 1-3 such as: !newsletter Denmark 2"}
}

func GetHelpMealplanMessageArray() []string {
	return []string{"Helper message for Mealplan :ramen:",
		"`",
		"**",
		"run !mealplan to get a breakfeast -> dinner -> snack mealplan for today!"}
}

func GetHelpSteamdealsMessageArray() []string {
	return []string{"Helper message for Steamdeals :video_game:",
		"`",
		"**",
		"call !steamdeals to get a list of 4 current deals on the steam store",
		"if you can less you can specify an amount: !steamdeals 2"}
}
