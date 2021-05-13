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
