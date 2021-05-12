package caching

import (
	"projectGroup23/structs"
	"time"

	"github.com/josemiguelmelo/gocacheable"
	"github.com/josemiguelmelo/gocacheable/providers/bigcache"
)

var cm = gocacheable.NewCacheableManager("manager_id")
var NewsCache structs.NewsLetters
var MealPlanCache structs.MealPlan

func AddCacheModule(name string) {
	cm.AddModule(name, &bigcache.BigCacheProvider{
		Lifetime: 100,
	})
}

func CacheNews(dur time.Duration) {
	// Caching the result of the GetWeatherForecast function
	err := cm.Cacheable("news", "newsletter", func() (interface{}, error) {
		res := NewsletterTest()
		return res, nil
	}, &NewsCache, dur)
	if err != nil {
		panic(err)
	}
}

func CacheMealplan(dur time.Duration) {
	// Caching the result of the GetWeatherForecast function
	err := cm.Cacheable("meals", "mealplan", func() (interface{}, error) {
		res := MealPlannerTest()
		return res, nil
	}, &MealPlanCache, dur)
	if err != nil {
		panic(err)
	}
}
