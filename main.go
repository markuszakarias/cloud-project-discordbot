package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

var server = "vmdata.database.windows.net"
var port = 1433
var user = "eriksen"
var password = "Tanzania1994!"
var database = "VM_Data"

var db *sql.DB

var err error

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

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
	fmt.Printf("Connected!\n")

	// Read employees
	count, err := getTodos()
	if err != nil {
		log.Fatal("Error reading from DB: ", err.Error())
	}
	fmt.Printf("Read %d row(s) successfully.\n", count)
}

func getTodos() (int, error) {
	ctx := context.Background()

	// Check database connection
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	// The sql query to be executed
	tsql := "SELECT Id, Userid, Title, Category, State FROM [dbo].[todo];"

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return -1, err
	}

	// Wait until function ends before
	defer rows.Close()

	var count int

	// Iterate through the result set.
	for rows.Next() {
		var title, userid, category, state string
		var id int

		// Get values from row.
		err := rows.Scan(&id, &userid, &title, &category, &state)
		if err != nil {
			return -1, err
		}

		fmt.Printf("Id:%d Userid:%s Title: %s, Category: %s, State: %s\n", id, userid, title, category, state)
		count++
	}

	return count, nil
}

/* 	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
   		server, user, password, port, database)

   	var err error

   	// Opening a database connection with tcp
   	db, err := sql.Open("sqlserver", connString)
   	if err != nil {
   		log.Fatal(err)
   	}

   	ctx := context.Background()
   	err = db.PingContext(ctx)
   	if err != nil {
   		log.Fatal(err.Error())
   	}
   	fmt.Printf("Connected!\n")

   	count, err := getTodos()
   	if err != nil {
   		log.Fatal("Error reading Employees: ", err.Error())
   	}
   	fmt.Printf("Read %d row(s) successfully.\n", count) */

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
// db.SetMaxIdleConns(100)
// Wait until the end of the function to close the connection
// defer db.Close()

// Creating an instance of the struct
//db_server := &utils.DB_server{DB: db}
//appHandler := &utils.RegexHandler{}

//fmt.Println(db_server)
//fmt.Println(appHandler)

//fmt.Println("Listening on port 8080")
//http.ListenAndServe(":8080", appHandler)

// Opening a database connection with tcp
