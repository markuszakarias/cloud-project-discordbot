// Package that contains our local caching. It adds a cache manager that handles all cache modules across the entire program.
// This package also manages the caching of each individual document in firestore.
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
		res, err := handlers.SteamDealsMainHandler(command)
		return res, err
	}, &DealsCache, dur)
	return err
}

// CacheForecasts - Caches the function return value and sets a timer for when the cache is dirty
func CacheForecasts(location string, dur time.Duration) error {
	// Cacheable adds cache to the function passed as parameter
	err := CM.Cacheable("cache", location, func() (interface{}, error) {
		res, err := handlers.WeatherForecastMainHandler(location)
		return res, err
	}, &ForecastsCache, dur)
	return err
}

// InitCache - Initializes the cache modules to be used in the bot
func InitCache() error {
	// Initializes BigCache cache
	err := AddCacheModule("cache")
	if err != nil {
		return err
	}
	return err
}

// CacheMeals - Caches the function return value and sets a timer for when the cache is dirty
func CacheMeals(dur time.Duration) error {
	// Cacheable adds cache to the function passed as parameter
	err := CM.Cacheable("cache", "meals", func() (interface{}, error) {
		res, err := handlers.MealPlanMainHandler()
		return res, err
	}, &MealsCache, dur)
	return err
}

// CacheNews - Caches the function return value and sets a timer for when the cache is dirty
func CacheNews(location string, dur time.Duration) error {
	// Cacheable adds cache to the function passed as parameter
	err := CM.Cacheable("cache", location, func() (interface{}, error) {
		res, err := handlers.NewsLetterMainHandler(location)
		return res, err
	}, &NewsCache, dur)
	return err
}
