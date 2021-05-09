package firebase

import (
	"context"
	"fmt"
	"log"
	"projectGroup23/handlers"
	"projectGroup23/structs"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Joke struct {
	text string
}

var ctx context.Context
var client *firestore.Client

// initialize firebase/firestore
func InitFirebase() {
	ctx = context.Background()
	opt := option.WithCredentialsFile("./firebase/firebasePrivateKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
}

func GetAllJokes() []string {
	iter := client.Collection("jokes").Documents(ctx)
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
		var test string = doc.Data()["text"].(string)
		allJokes = append(allJokes, test)
	}
	return allJokes
}

func CreateJoke(userId string, jokeText string) error {
	_, _, err := client.Collection("jokes").Add(ctx,
		map[string]interface{}{
			"createdBy": userId,
			"text":      jokeText,
		})
	return err
}

func GetAllJokesByUserId(userId string) []string {
	//iter := client.Collection("jokes").Documents(ctx)
	iter := client.Collection("jokes").Where("createdBy", "==", userId).Documents(ctx)
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

func getAllWebhooks() ([]structs.CloudWebhook, error) {
	iter := client.Collection("cloudwebhook").Documents(ctx)
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
			Id:                   doc.Ref.ID,
			UserId:               doc.Data()["UserId"].(string),
			CloudPercentages:     doc.Data()["CloudPercentages"].(int64),
			HasBeenNotifiedToday: doc.Data()["HasBeenNotifiedToday"].(bool),
		}
		allWebhooks = append(allWebhooks, webhookPlaceholder)
	}
	return allWebhooks, nil
}

func DeleteWebhook(userId string) error {
	_, err := client.Collection("cloudwebhook").Doc(userId).Delete(ctx)
	return err
}

// if user already has a weather webhook, it will be updated!
func CreateWeatherWebhook(userId string, cloudPercentages int64) error {
	_, err := client.Collection("cloudwebhook").Doc(userId).Set(ctx, map[string]interface{}{
		"Id":                   "",
		"UserId":               userId,
		"CloudPercentages":     cloudPercentages,
		"HasBeenNotifiedToday": false,
	}, firestore.MergeAll)
	return err
}

func updateWeatherWebhook(userId string, webhookData map[string]interface{}) error {
	_, err := client.Collection("cloudwebhook").Doc(userId).Set(ctx, webhookData, firestore.MergeAll)
	return err
}

func WebhookRoutine(s *discordgo.Session) {
	webhooks, err := getAllWebhooks()
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(webhooks); i++ {
		wf := handlers.GetWeatherForecast(1)
		currentCloud := wf.Forecasts[0].Clouds
		if float64(webhooks[i].CloudPercentages) >= currentCloud && !webhooks[i].HasBeenNotifiedToday { // if less cloud than notification setting and has not been notified today
			userChannel, _ := s.UserChannelCreate(webhooks[i].UserId)
			message := "tomorrow it will be " + fmt.Sprintf("%.f", currentCloud) + " percent cloud!"
			s.ChannelMessageSend(userChannel.ID, message)
			webhookData := map[string]interface{}{
				"HasBeenNotifiedToday": true,
			}
			err := updateWeatherWebhook(webhooks[i].UserId, webhookData) // updates the webbook so it can't notify again today
			if err != nil {
				log.Fatalln("An error has occurred: %s", err)
			}
		}
	}
	time.Sleep(time.Duration(900) * time.Second) // runs every 15 minuts
}
