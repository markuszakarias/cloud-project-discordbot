package caching

import (
	"github.com/josemiguelmelo/gocacheable"
	"github.com/josemiguelmelo/gocacheable/providers/bigcache"
	"projectGroup23/structs"
	"time"
)

var cm = gocacheable.NewCacheableManager("manager_id")
var NewsCache structs.NewsLetters

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
