package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

	// Opening a database connection with tcp
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(db)

	// "root:@tcp(127.0.0.1:3306)/todo_app"

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(100)
	// Wait until the end of the function to close the connection
	defer db.Close()

	// Creating an instance of the struct
	db_server := &utils.DB_server{DB: db}
	appHandler := &utils.RegexHandler{}

	fmt.Println(db_server)
	fmt.Println(appHandler)

	// Handlers for the endpoints. Mainly used for testing the output of the database interactions.
	// Uses regex for URL
	appHandler.Handler("/app/$", "GET", db_server.TodoAll)
	/*appHandler.Handler("/app/$", "POST", db_server.TodoCreate)
	appHandler.Handler("/app/[0-9]+$", "GET", db_server.TodoObject)
	appHandler.Handler("/app/[0-9]+$", "PUT", db_server.TodoUpdate)
	appHandler.Handler("/app/[0-9]+$", "DELETE", db_server.TodoDeleteObject)
	*/
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", appHandler)

}
