package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello This is a simple thing")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	/* 	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyDyKyvfKQuu90-cieJZ2vPiK9u6n7DLWFE"))
	   	if err != nil {
	   		log.Fatalf("fatal error: %s", err)
	   	}

	   	r := &maps.DirectionsRequest{
	   		Origin:      "Oslo",
	   		Destination: "Trondheim",
	   	}

	   	route, _, err := c.Directions(context.Background(), r)
	   	if err != nil {
	   		log.Fatalf("fatal error: %s", err)
	   	}

	   	pretty.Println(route)
	*/

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", simpleHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
