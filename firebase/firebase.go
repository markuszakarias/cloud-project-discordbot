package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
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
		var test string = doc.Data()["text"].(string)
		allJokes = append(allJokes, test)
	}
	return allJokes
}
