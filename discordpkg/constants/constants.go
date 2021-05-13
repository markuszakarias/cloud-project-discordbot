package constants

func GetWeatherStringArray() []string {
	return []string{":calendar: ", ":map: ", " - ", ":cloud: ", ":dash: ",
		" m/s", " Celsius", "Probability of precipitation: ", ":cloud_rain: ", ":cloud_snow: ",
		"Temperature:", ":city_sunrise: ", ":cityscape: ", ":city_dusk: ", ":night_with_stars: ", "%"}
}

func GetSteamStringArray() []string {
	return []string{"title: ", "DealID: ", "NormalPrice: ", "SalePrice: ", "Savings: ",
		"MetacriticScore: ", "SteamRatingText: ", "SteamRatingPercent: ", "SteamRatingCount: "}
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
	return []string{"Helper message for Weather :partly_sunny:"}
}

func GetHelpNewsletterMessageArray() []string {
	return []string{"Helper message for Newsletter :earth_africa:"}
}

func GetHelpMealplanMessageArray() []string {
	return []string{"Helper message for Mealplan :ramen:"}
}

func GetHelpSteamdealsMessageArray() []string {
	return []string{"Helper message for Steamdeals :video_game:"}
}
