package structs

// NewsLetter contains the information for one single news headline from newsapi.
type NewsLetter struct {
	Author         string `json:"author"`
	Date_published string `json:"date_published"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Url_to_story   string `json:"url_to_story"`
}

// NewsLetters prepares a array with 5 newsletters.
type NewsLetters struct {
	Newsletters [5]NewsLetter `json:"newsletters"`
}

// Meal holds information of one meal from the external spoonacular api.
type Meal struct {
	Title          string `json:"title"`
	ReadyInMinutes string `json:"ready in minutes"`
	Url            string `json:"url"`
}

// mealPlan is the plan that is given through the discord bot to the user when asked for a meal plan.
type MealPlan struct {
	MealMessage string  `json:"mealMessage"`
	Meals       [3]Meal `json:"meals"`
	Nutrients   struct {
		Calories      float64 `json:"calories"`
		Protein       float64 `json:"protein"` // should be floating point numbers!
		Fat           float64 `json:"fat"`
		CarboHydrates float64 `json:"carbohydrates"`
	}
}

type IPPosition struct {
	City	string `json:"city"`
	Country string `json:"country_name"`
}

type WeatherForecast struct {

}
