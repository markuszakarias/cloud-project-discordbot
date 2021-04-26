package main

import (
	"context"
	"log"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

func main() {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyDyKyvfKQuu90-cieJZ2vPiK9u6n7DLWFE"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.DirectionsRequest{
		Origin:      "Los Angeles",
		Destination: "New York",
	}

	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Println(route)
}
