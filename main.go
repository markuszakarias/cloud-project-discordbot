package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type todo_struct struct {
	id       int64
	userid   string
	title    string
	category string
	state    string
}

var ctx context.Context
var client *firestore.Client

var collection = "todolist"

func InitFirebase() {
	ctx = context.Background()

	opt := option.WithCredentialsFile("./assets/assignment-2-db-firebase-adminsdk-ij2bm-1dbf677c2d.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		fmt.Println("This error happened when connecting to firestore.")
		log.Fatalln(err)
	}

	// Used a the main todo object
	var todoObject todo_struct
	todoObject.userid = "001"
	todoObject.title = "Second Task"
	todoObject.category = "Updated Category"
	todoObject.state = "inactive"

	// Used for updating a todo object
	var updateObject todo_struct
	updateObject.userid = "001"
	updateObject.title = "Second Task"
	updateObject.category = "Updated Category"
	updateObject.state = "active"

	/* getlist := getTodoAll()
	fmt.Println(getlist) */

	/* getTodo := getTodoObject(todoObject)
	fmt.Println(getTodo) */

	/* createTodo := createTodoObject(todoObject)
	fmt.Println(createTodo) */

	/* err = deleteTodoObject(todoObject)
	if err != nil {
		fmt.Println(err)
	} */

	/* err = updateTodoObject(todoObject, updateObject)
	if err != nil {
		fmt.Println(err)
	} */

}

func getTodoAll() []todo_struct {
	iter := client.Collection(collection).Documents(ctx)

	var objects []todo_struct
	var todoObject todo_struct

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err) //TODO error handle this
		}

		todoObject.userid = doc.Data()["Userid"].(string)
		todoObject.title = doc.Data()["Title"].(string)
		todoObject.category = doc.Data()["Category"].(string)
		todoObject.state = doc.Data()["State"].(string)
		objects = append(objects, todoObject)
	}
	return objects
}

func getTodoObject(todoObject todo_struct) todo_struct {
	data, err := client.Collection(collection).Doc(todoObject.userid).Get(ctx)
	if err != nil {
		fmt.Println(err) //TODO error handle this
	}

	if data.Data() == nil {
		fmt.Println("data is nil") //TODO error handle this
	}
	todo := data.Data()
	todo["Userid"] = todoObject.userid

	fmt.Println(todo)

	return todoObject
}

func createTodoObject(todoObject todo_struct) todo_struct {
	_, _, err := client.Collection(collection).Add(ctx,
		map[string]interface{}{
			"Userid":   todoObject.userid,
			"Title":    todoObject.title,
			"Category": todoObject.category,
			"State":    todoObject.state,
		})
	if err != nil {
		return todoObject
	}

	return todoObject
}

func updateTodoObject(todoObject todo_struct, updateObject todo_struct) error {
	query := client.Collection(collection).Where("Userid", "==", todoObject.userid).Where("Title", "==", todoObject.title).Documents(ctx)

	for {
		doc, err := query.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		todo_id := doc.Ref.ID
		client.Collection(collection).Doc(todo_id).Set(ctx, map[string]interface{}{
			"Title":    updateObject.title,
			"Category": updateObject.category,
			"State":    updateObject.state,
		}, firestore.MergeAll)
	}

	return nil
}

func deleteTodoObject(todoObject todo_struct) error {
	query := client.Collection(collection).Where("Userid", "==", todoObject.userid).Where("Title", "==", todoObject.title).Documents(ctx)

	for {
		doc, err := query.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		todo_id := doc.Ref.ID
		res, err := client.Collection(collection).Doc(todo_id).Delete(ctx)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}

	return nil
}

func main() {
	InitFirebase()
}
