package caching

import (
	"projectGroup23/structs"
	"time"

	"github.com/josemiguelmelo/gocacheable"
	"github.com/josemiguelmelo/gocacheable/providers/bigcache"
)

var cm = gocacheable.NewCacheableManager("manager_id")
var DealsCache structs.Deals
var ForecastsCache structs.WeatherForecasts
var MealsCache structs.MealPlan
var NewsCache structs.NewsLetters



func AddCacheModule(name string) {
	cm.AddModule(name, &bigcache.BigCacheProvider{
		Lifetime: 100,
	})
}

func CacheDeals(command string, dur time.Duration) {
	// Caching the result of the GetWeatherForecast function
	err := cm.Cacheable("deals", "deals", func() (interface{}, error) {
		res := SteamdealsTest(command)
		return res, nil
	}, &DealsCache, dur)
	if err != nil {
		panic(err)
	}
}

func CacheForecasts(dur time.Duration) {
	// Caching the result of the GetWeatherForecast function
	err := cm.Cacheable("forecasts", "forecasts", func() (interface{}, error) {
		res := WeatherForecastTest()
		return res, nil
	}, &ForecastsCache, dur)
	if err != nil {
		panic(err)
	}
}

func CacheMeals(dur time.Duration) {
	// Caching the result of the GetWeatherForecast function
	err := cm.Cacheable("meals", "meals", func() (interface{}, error) {
		res := MealPlannerTest()
		return res, nil
	}, &MealsCache, dur)
	if err != nil {
		panic(err)
	}
}

func CacheNews(dur time.Duration) {
	// Caching the result of the GetWeatherForecast function
	err := cm.Cacheable("news", "news", func() (interface{}, error) {
		res := NewsletterTest()
		return res, nil
	}, &NewsCache, dur)
	if err != nil {
		panic(err)
	}
}