package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"projectGroup23/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := ":" + port
	fmt.Print("Listening on ", address, "...\n")

	http.HandleFunc("/newsletter", handlers.GetDailyNewsLetter)
	http.HandleFunc("/mealplan", handlers.GetDailyMealPlan)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
