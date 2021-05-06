package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

type todo_struct struct {
	id       int64
	userid   string
	title    string
	category string
	state    string
}

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

	// Use this stuct to create/display/delete/update one todo object
	var todoObject todo_struct
	todoObject.userid = "002"
	todoObject.title = "Second Task"
	todoObject.category = "Second Category"
	todoObject.state = "active"

	// Read all todo
	err = getTodoAll()
	if err != nil {
		log.Fatal("Error reading all todo objects: ", err.Error())
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

func updateTodoObject(todoObject todo_struct) error {
	ctx := context.Background()

	if db == nil {
		err = errors.New("create todo object: db is null")
	}

	// Check db connection
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := "UPDATE [dbo].[todo] SET Category = @Category WHERE Userid = @Userid;"

	// Execute non-query with named parameters
	res, err := db.ExecContext(ctx, tsql,
		sql.Named("Userid", todoObject.userid),
		sql.Named("Category", todoObject.category))

	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func deleteTodoObject(todoObject todo_struct) error {
	ctx := context.Background()

	if db == nil {
		err = errors.New("create todo object: db is null")
	}

	// Check db connection
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := "DELETE FROM [dbo].[todo] WHERE Userid = @Userid;"

	// Execute non-query with named parameters
	res, err := db.ExecContext(ctx, tsql,
		sql.Named("Userid", todoObject.userid))

	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func createTodoObject(todoObject todo_struct) error {
	ctx := context.Background()

	if db == nil {
		err = errors.New("create todo object: db is null")
	}

	// Check db connection
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := `INSERT INTO [dbo].[todo] (Userid, Title, Category, State) VALUES (@Userid, @Title, @Category, @State);
			select isNull(SCOPE_IDENTITY(), -1);`

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx,
		sql.Named("Userid", todoObject.userid),
		sql.Named("Title", todoObject.title),
		sql.Named("Category", todoObject.category),
		sql.Named("State", todoObject.state))

	fmt.Println("Before scan")

	err = row.Scan(&todoObject.id)
	if err != nil {
		return err
	}

	fmt.Println("New ID was created:", todoObject.id)

	return nil
}

func getTodoAll() error {
	ctx := context.Background()

	// Check database connection
	err := db.PingContext(ctx)
	if err != nil {
		return err
	}

	// The sql query to be executed
	tsql := "SELECT Id, Userid, Title, Category, State FROM [dbo].[todo];"

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return err
	}

	// Wait until function ends before
	defer rows.Close()

	fmt.Println("Print all todo data:")

	// Iterate through the result set.
	for rows.Next() {
		var title, userid, category, state string
		var id int

		// Get values from row.
		err := rows.Scan(&id, &userid, &title, &category, &state)
		if err != nil {
			return err
		}

		fmt.Printf("Id:%d Userid:%s Title: %s, Category: %s, State: %s\n", id, userid, title, category, state)
	}

	return nil
}

func getTodoObject(todoObject todo_struct) error {
	ctx := context.Background()

	// Check database connection
	err := db.PingContext(ctx)
	if err != nil {
		return err
	}

	// The sql query to be executed
	tsql := "SELECT Id, Userid, Title, Category, State FROM [dbo].[todo] WHERE Userid = @Userid;"

	// Execute query
	rows, err := db.QueryContext(ctx, tsql, sql.Named("UserId", todoObject.userid))
	if err != nil {
		return err
	}

	// Wait until function ends before
	defer rows.Close()

	fmt.Println("Print one todo object:")

	// Iterate through the result set.
	for rows.Next() {

		// Get values from row.
		err := rows.Scan(&todoObject.id,
			&todoObject.userid,
			&todoObject.title,
			&todoObject.category,
			&todoObject.state)

		if err != nil {
			return err
		}

		fmt.Printf("Id:%d Userid:%s Title: %s, Category: %s, State: %s\n",
			todoObject.id, todoObject.userid, todoObject.title, todoObject.category, todoObject.state)
	}

	return nil
}
