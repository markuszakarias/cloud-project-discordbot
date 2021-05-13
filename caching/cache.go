package caching

import (
	"projectGroup23/handlers"
	"projectGroup23/structs"
	"time"

	"github.com/josemiguelmelo/gocacheable"
	"github.com/josemiguelmelo/gocacheable/providers/bigcache"
)

// CM - Manages all cache modules across the entire program
var CM = gocacheable.NewCacheableManager("manager_id")

// Global variables for storing function returns from cache
var DealsCache structs.Deals
var ForecastsCache structs.WeatherForecasts
var MealsCache structs.MealPlan
var NewsCache structs.NewsLetters


// AddCacheModule - Initializes a cache module with BigCache and HardMaxCacheSize: 8192 MB
func AddCacheModule(name string) error {
	// We use BigCache since reads are lock-free and its best for read-only functionality
	err := CM.AddModule(name, &bigcache.BigCacheProvider{
		Lifetime: 100,
	})
	if err != nil {
		return err
	}
	return err
}

// CacheDeals - Caches the function return value and sets a timer for when the cache is dirty
func CacheDeals(command string, dur time.Duration) error {
	// Cacheable adds cache to the function passed as parameter
	err := CM.Cacheable("cache", "deals", func() (interface{}, error) {
		res := handlers.SteamDealsMainHandler(command)
		return res, nil
	}, &DealsCache, dur)
	if err != nil {
		return err
	}
	return err
}

// CacheForecasts - Caches the function return value and sets a timer for when the cache is dirty
func CacheForecasts(apikey string, dur time.Duration) error {
	// Cacheable adds cache to the function passed as parameter
	err := CM.Cacheable("cache", "forecasts", func() (interface{}, error) {
		res := handlers.WeatherForecastMainHandler(apikey)
		return res, nil
	}, &ForecastsCache, dur)
	if err != nil {
		return err
	}
	return err
}

// CacheMeals - Caches the function return value and sets a timer for when the cache is dirty
func CacheMeals(apikey string, dur time.Duration) error {
	// Cacheable adds cache to the function passed as parameter
	err := CM.Cacheable("cache", "meals", func() (interface{}, error) {
		res := handlers.MealPlanMainHandler(apikey)
		return res, nil
	}, &MealsCache, dur)
	if err != nil {
		return err
	}
	return err
}

// CacheNews - Caches the function return value and sets a timer for when the cache is dirty
func CacheNews(apikey string, dur time.Duration) error {
	// Cacheable adds cache to the function passed as parameter
	err := CM.Cacheable("cache", "news", func() (interface{}, error) {
		res := handlers.NewsLetterMainHandler(apikey)
		return res, nil
	}, &NewsCache, dur)
	if err != nil {
		return err
	}
	return err
}
