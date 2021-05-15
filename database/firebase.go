package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"projectGroup23/structs"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Joke struct {
	text string
}

var Ctx context.Context
var Client *firestore.Client

// struct used for the stored data
var StoredMealPlan structs.StoredMealPlan

// struct used for the cached data
var StoredWeatherForecast structs.StoredWeatherForecast

// struct used for the cached data
var StoredSteamDeals structs.StoredSteamDeals

// struct used for the cached data
var StoredNewsLetter structs.StoredNewsLetter

// initialize firebase/firestore
func InitFirebase() {
	Ctx = context.Background()
	opt := option.WithCredentialsFile("./database/firebasePrivateKey.json")
	app, err := firebase.NewApp(Ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	Client, err = app.Firestore(Ctx)
}

func GetRandomJoke() (string, error) {
	iter := Client.Collection("jokes").Documents(Ctx)
	var allJokes []string
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
		//joke := doc.Data()["text"]
		var test string = doc.Data()["text"].(string)
		allJokes = append(allJokes, test)
	}
	if len(allJokes) == 0 { // if no jokes in firebase
		return "", errors.New("No jokes in database. Create one with the '!joke create' command")
	}

	rand.Seed(time.Now().Unix())                     // pseudo random
	randomJoke := allJokes[rand.Intn(len(allJokes))] // takes a random joke from allJokes
	return randomJoke, nil                           // returns the random joke
}

func CreateJoke(userId string, jokeText string) error {
	_, _, err := Client.Collection("jokes").Add(Ctx,
		map[string]interface{}{
			"createdBy": userId,
			"text":      jokeText,
		})
	return err
}

func GetAllJokesByUserId(userId string) []string {
	//iter := client.Collection("jokes").Documents(ctx)
	iter := Client.Collection("jokes").Where("createdBy", "==", userId).Documents(Ctx)
	var allJokes []string
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return allJokes
		}
		//joke := doc.Data()["text"]
		var joketext string = doc.Data()["text"].(string)
		allJokes = append(allJokes, joketext)
	}
	return allJokes
}

func GetAllWebhooks() ([]structs.CloudWebhook, error) {
	iter := Client.Collection("cloudwebhook").Documents(Ctx)
	var allWebhooks []structs.CloudWebhook
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return allWebhooks, err
		}
		//joke := doc.Data()["text"]
		webhookPlaceholder := structs.CloudWebhook{
			UserId: doc.Data()["UserId"].(string),
			City:   doc.Data()["City"].(string),
		}
		allWebhooks = append(allWebhooks, webhookPlaceholder)
	}
	return allWebhooks, nil
}

func DeleteWebhook(userId string) error {
	_, err := Client.Collection("cloudwebhook").Doc(userId).Delete(Ctx)
	return err
}

// if user already has a weather webhook, it will be updated!
func CreateWeatherWebhook(userId string, city string) error {
	_, err := Client.Collection("cloudwebhook").Doc(userId).Set(Ctx, map[string]interface{}{
		"Id":     "",
		"UserId": userId,
		"City":   city,
		// "CloudPercentages": cloudPercentages,
		// "LastDateNotified": time.Now().AddDate(0, 0, -1), // sets the dat before since it has not been notified today yet.
	}, firestore.MergeAll)
	return err
}

func updateWeatherWebhook(userId string, webhookData map[string]interface{}) error {
	_, err := Client.Collection("cloudwebhook").Doc(userId).Set(Ctx, webhookData, firestore.MergeAll)
	return err
}

func CheckWeatherForecastsOnFirestore(location string) (structs.StoredWeatherForecast, error) {
	iter := Client.Collection("cached_resp").Documents(Ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if doc.Data()["IPLocation"] == location {
			doc.DataTo(&StoredWeatherForecast)
			StoredWeatherForecast.FirestoreID = doc.Ref.ID // matching the firestore ID with the one stored
			return StoredWeatherForecast, nil
		}
	}

	return structs.StoredWeatherForecast{}, nil
}

// GetCachedNewsLetterFromFirestore - global function that runs at startup
// gets all the cached data from firestore
func GetStoredFromFirestore() {
	iter := Client.Collection("cached_resp").Documents(Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if doc.Data()["WeatherForecasts"] != nil {
			doc.DataTo(&StoredWeatherForecast)
			StoredWeatherForecast.FirestoreID = doc.Ref.ID // matching the firestore ID with the one stored
		}

		if doc.Data()["NewsLetters"] != nil {
			doc.DataTo(&StoredNewsLetter)
			StoredNewsLetter.FirestoreID = doc.Ref.ID // matching the firestore ID with the one stored
		}

		if doc.Data()["SteamDeals"] != nil {
			doc.DataTo(&StoredSteamDeals)
			StoredSteamDeals.FirestoreID = doc.Ref.ID // matching the firestore ID with the one stored
		}

		if doc.Data()["MealPlan"] != nil {
			doc.DataTo(&StoredMealPlan)
			StoredMealPlan.FirestoreID = doc.Ref.ID // matching the firestore ID with the one stored
		}
	}
}

// DeleteObjectFromFirestore - deletes an object in firestore based on firestore ID
func DeleteObjectFromFirestore(firestoreID string) {
	_, err := Client.Collection("cached_resp").Doc(firestoreID).Delete(Ctx)
	if err != nil {
		fmt.Println(err)
	}
}

// UpdateTimeFirestore - updates the object in firestore
func UpdateTimeFirestore(firestoreID string, storeTime time.Time, storeRefresh float64) {
	_, err := Client.Collection("cached_resp").Doc(firestoreID).Update(Ctx, []firestore.Update{
		{
			Path:  "StoreTime", // matching specific field in firestore object
			Value: storeTime,
		},
		{
			Path:  "StoreRefresh", // matching specific field in firestore object
			Value: storeRefresh,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

// saveNewsLetterToFirestore - saves an object to firestore
func SaveSteamDealsToFirestore(stored *structs.StoredSteamDeals) error {
	doc, _, err := Client.Collection("cached_resp").Add(Ctx, *stored)
	stored.FirestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(stored.FirestoreID) // confirming the storage of document ID
	return err
}

// saveNewsLetterToFirestore - saves an object to firestore
func SaveMealPlannerToFirestore(stored *structs.StoredMealPlan) error {
	doc, _, err := Client.Collection("cached_resp").Add(Ctx, *stored)
	stored.FirestoreID = doc.ID // storing firestore ID for later use
	return err
}

// saveNewsLetterToFirestore - saves an object to firestore
func SaveNewsLetterToFirestore(stored *structs.StoredNewsLetter) error {
	doc, _, err := Client.Collection("cached_resp").Add(Ctx, *stored)
	stored.FirestoreID = doc.ID // storing firestore ID for later use
	return err
}

// saveNewsLetterToFirestore - saves an object to firestore
func SaveWeatherForecastToFirestore(stored *structs.StoredWeatherForecast) error {
	_, _, err := Client.Collection("cached_resp").Add(Ctx, *stored)
	//stored.FirestoreID = doc.ID // storing firestore ID for later use
	return err
}
