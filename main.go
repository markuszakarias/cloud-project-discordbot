package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"projectGroup23/structs"
	"projectGroup23/utils"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

var server = "vmdata.database.windows.net"
var port = 1433
var user = "eriksen"
var password = "Tanzania1994!"
var database = "VM_Data"

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	// Create connection pool
	utils.Db, utils.Err = sql.Open("sqlserver", connString)
	if utils.Err != nil {
		log.Fatal("Error creating connection pool: ", utils.Err.Error())
	}
	ctx := context.Background()
	utils.Err = utils.Db.PingContext(ctx)
	if utils.Err != nil {
		log.Fatal(utils.Err.Error())
	}

	// Use this stuct to create/display/delete/update one todo object
	var todoObject structs.Todo_struct
	todoObject.Userid = "002"
	todoObject.Title = "Second Task"
	todoObject.Category = "Second Category"
	todoObject.State = "active"

	// Read all todo
	utils.Err = utils.GetTodoAll()
	if utils.Err != nil {
		log.Fatal("Error reading all todo objects: ", utils.Err.Error())
	}

	// Read one todo
	/* err = getTodoObject("abcdefgh12345678")
	if err != nil {
		log.Fatal("Error reading todo object: ", err.Error())
	} */

	// Create todo object
	/*err = createTodoObject(todoObject)
	if err != nil {
		log.Fatal("Error creating todo object: ", err.Error())
	} */

	// Delete todo object
	/* 	err = deleteTodoObject("abcdefgh12345678")
	   	if err != nil {
	   		log.Fatal("Error reading all todo objects: ", err.Error())
	   	} */

	// Update todo object
	/* 	err = updateTodoObject(todoObject)
	   	if err != nil {
	   		log.Fatal("Error updating todo object: ", err.Error())
	   	} */
}
