// Package database - contains our handling of the Azure SQL and Firestore connection. It initiates the connection and has the functionality to
// to view, create, update and delete todo tasks from azure sql. And also view and create jokes from Firestore. It also uses the environment
// variables for the database credentials. It also contains functionality for handling webhooks and stored struct object regarding the weather,
// steamdeals, newsletter and mealplan API data.
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"projectGroup23/structs"
	"projectGroup23/utils"
)

// variable for the database connection
var db *sql.DB

// global package error variable
var err error

// variables to establish azure sql database connection
var server = utils.EnvVar("DB_SERVER")
var port = utils.EnvVar("DB_PORT")
var user = utils.EnvVar("DB_USER")
var password = utils.EnvVar("DB_PASSWORD")
var azuredb = utils.EnvVar("DB")

// InitSQL - Used in main to initialize the azure sql database connection
func InitSQL() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, azuredb)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// UpdateTodoObject - Updates a todo object from the ID with a new description
func UpdateTodoObject(sqlId int, description string) error {

	ctx := context.Background()
	if db == nil {
		err = errors.New("create todo object: db is null")
	}

	// Check db connection
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	// The sql query
	tsql := "UPDATE [dbo].[todo] SET Description = @Description WHERE Id = @Id;"

	// Execute non-query with named parameters
	_, err := db.ExecContext(ctx, tsql,
		sql.Named("Id", sqlId),
		sql.Named("Description", description))

	if err != nil {
		return err
	}

	return nil
}

// UpdateTodoObjectStatus - Updates a todo object from the ID with a status
func UpdateTodoObjectStatus(sqlId int, status string) error {

	ctx := context.Background()
	if db == nil {
		return errors.New("create todo object: db is null")
	}

	// Check db connection
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	// The sql query
	tsql := "UPDATE [dbo].[todo] SET State = @State WHERE Id = @Id;"

	// Execute non-query with named parameters
	_, err := db.ExecContext(ctx, tsql,
		sql.Named("Id", sqlId),
		sql.Named("State", status))

	if err != nil {
		return err
	}

	return nil
}

// DeleteTodoObject - Deletes a todo object matching the ID
func DeleteTodoObject(sqlId int) error {

	ctx := context.Background()
	if db == nil {
		return errors.New("create todo object: db is null")
	}

	// Check db connection
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	// The sql query
	tsql := "DELETE FROM [dbo].[todo] WHERE Id = @Id;"

	// Execute non-query with named parameters
	_, err := db.ExecContext(ctx, tsql,
		sql.Named("Id", sqlId))

	if err != nil {
		return err
	}

	return nil
}

// CreateTodoObject - Creates a todo object with new id and description from the user
func CreateTodoObject(todoObject structs.TodoStruct) error {

	ctx := context.Background()
	if db == nil {
		return errors.New("create todo object: db is null")
	}

	// Check db connection
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	// The sql query
	tsql := `INSERT INTO [dbo].[todo] (Userid, Description, State) VALUES (@Userid, @Description, @State);
			select isNull(SCOPE_IDENTITY(), -1);`

	// prepares the sql insertion
	stmt, err := db.Prepare(tsql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	// insert the data into a row
	row := stmt.QueryRowContext(ctx,
		sql.Named("Userid", todoObject.Userid),
		sql.Named("Description", todoObject.Description),
		sql.Named("State", todoObject.State))

	err = row.Scan(&todoObject.Id)
	if err != nil {
		return err
	}

	return nil
}

// GetTodoObject - Displayes the todo tasks of the user running the command
func GetTodoObject(userid string) ([]structs.TodoStruct, error) {

	// list of todo structs to insert todo object(s)
	var allTodos []structs.TodoStruct

	ctx := context.Background()

	// Check database connection
	err := db.PingContext(ctx)
	if err != nil {
		return allTodos, err
	}

	// The sql query
	tsql := "SELECT Id, Userid, Description, State FROM [dbo].[todo] WHERE Userid = @Userid;"

	// Execute query
	rows, err := db.QueryContext(ctx, tsql, sql.Named("Userid", userid))
	if err != nil {
		return allTodos, err
	}

	// Wait until function ends before
	defer rows.Close()

	// Iterate through the result set.
	for rows.Next() {
		var userid, description, state string
		var id int

		// Get values from row.
		err := rows.Scan(&id, &userid, &description, &state)
		if err != nil {
			return allTodos, err
		}

		todoItem := structs.TodoStruct{
			Id:          id,
			Userid:      userid,
			Description: description,
			State:       state,
		}

		allTodos = append(allTodos, todoItem)
	}

	return allTodos, nil
}

// ConvertIndexToId - helper function to convert the ID of a todo task, to the ID in azure sql
func ConvertIndexToId(i int, userid string) (int, error) {
	resp, err := GetTodoObject(userid)
	if err != nil {
		return 0, err
	}
	if len(resp) <= i {
		return 0, errors.New("id does not exist")
	}
	convID := resp[i].Id

	return convID, nil
}
