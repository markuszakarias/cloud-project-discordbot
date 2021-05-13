package structs

import "time"

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
	Newsletters []NewsLetter `json:"newsletters"`
}

// Meal holds information of one meal from the external spoonacular api.
type Meal struct {
	Title          string `json:"title"`
	ReadyInMinutes string `json:"ready in minutes"`
	Url            string `json:"url"`
}

// mealPlan is the plan that is given through the discord bot to the user when asked for a meal plan.
type MealPlan struct {
	MealMessage string `json:"mealMessage"`
	Meals       []Meal `json:"meals"`
	Nutrients   struct {
		Calories      float64 `json:"calories"`
		Protein       float64 `json:"protein"` // should be floating point numbers!
		Fat           float64 `json:"fat"`
		CarboHydrates float64 `json:"carbohydrates"`
	}
}

type IPLocation struct {
	City string `json:"city"`
}

type WeatherForecast struct {
	Date    string  `json:"date"`
	City    string  `json:"city"`
	Main    string  `json:"main"`
	Desc    string  `json:"description"`
	Morning float64 `json:"morning"`
	Day     float64 `json:"day"`
	Eve     float64 `json:"eve"`
	Night   float64 `json:"night"`
	Clouds  float64 `json:"clouds"`
	Wind    float64 `json:"wind"`
	POP     float64 `json:"pop"`
	Rain    float64 `json:"rain"`
	Snow    float64 `json:"snow"`
}

type WeatherForecasts struct {
	Forecasts []WeatherForecast `json:"forecasts"`
}

// Deal contains the information for what we have chosen to define as a deal.
type Deal struct {
	Title              string `json:"title"`
	DealID             string `json:"deal-id"`
	NormalPrice        string `json:"normal price"`
	SalePrice          string `json:"sale price"`
	Savings            string `json:"savings"`
	MetacriticScore    string `json:"meta critic score"`
	SteamRatingText    string `json:"steam rating"`
	SteamRatingPercent string `json:"steam rating percent"`
	SteamRatingCount   string `json:"steam rating amount"`
}

// Deals contains a slice of deal structs, used when presenting with multiple or 1 deal.
type Deals struct {
	Deals []Deal
}

// Struct todo_struct - struct for a todo object
type Todo_struct struct {
	Id          int
	Userid      string
	Description string
	State       string
}

type CloudWebhook struct {
	Id               string    `json:"id"`
	UserId           string    `json:"userid"`
	CloudPercentages int64     `json:"cloudpercentages"`
	LastDateNotified time.Time `json:"lastdatenotified"`
}

type StoredMealPlan struct {
	MealPlan     MealPlan
	StoreTime    time.Time
	StoreRefresh float64
	FirestoreID  string
}

// StoredNewsLetter - struct for a stored newsletter
type StoredNewsLetter struct {
	NewsLetters  NewsLetters
	StoreTime    time.Time
	StoreRefresh float64
	FirestoreID  string
}

// CachedNewsLetter - struct for a cached newsletter
type StoredSteamDeals struct {
	SteamDeals   Deals
	StoreTime    time.Time
	StoreRefresh float64
	FirestoreID  string
}

type StoredWeatherForecast struct {
	WeatherForecasts WeatherForecasts
	IPLocation       IPLocation
	StoreTime       time.Time
	StoreRefresh    float64
	FirestoreID      string
}
