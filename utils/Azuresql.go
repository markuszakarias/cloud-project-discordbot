package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"projectGroup23/structs"
)

var Db *sql.DB

var Err error

func UpdateTodoObject(todoObject structs.Todo_struct) error {
	ctx := context.Background()

	if Db == nil {
		Err = errors.New("create todo object: db is null")
	}

	// Check db connection
	Err = Db.PingContext(ctx)
	if Err != nil {
		return Err
	}

	tsql := "UPDATE [dbo].[todo] SET Category = @Category WHERE Userid = @Userid;"

	// Execute non-query with named parameters
	res, err := Db.ExecContext(ctx, tsql,
		sql.Named("Userid", todoObject.Userid),
		sql.Named("Category", todoObject.Category))

	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func DeleteTodoObject(todoObject structs.Todo_struct) error {
	ctx := context.Background()

	if Db == nil {
		Err = errors.New("create todo object: db is null")
	}

	// Check db connection
	Err = Db.PingContext(ctx)
	if Err != nil {
		return Err
	}

	tsql := "DELETE FROM [dbo].[todo] WHERE Userid = @Userid;"

	// Execute non-query with named parameters
	res, err := Db.ExecContext(ctx, tsql,
		sql.Named("Userid", todoObject.Userid))

	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func CreateTodoObject(todoObject structs.Todo_struct) error {
	ctx := context.Background()

	if Db == nil {
		Err = errors.New("create todo object: db is null")
	}

	// Check db connection
	Err = Db.PingContext(ctx)
	if Err != nil {
		return Err
	}

	tsql := `INSERT INTO [dbo].[todo] (Userid, Title, Category, State) VALUES (@Userid, @Title, @Category, @State);
			select isNull(SCOPE_IDENTITY(), -1);`

	stmt, err := Db.Prepare(tsql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx,
		sql.Named("Userid", todoObject.Userid),
		sql.Named("Title", todoObject.Title),
		sql.Named("Category", todoObject.Category),
		sql.Named("State", todoObject.State))

	err = row.Scan(&todoObject.Id)
	if err != nil {
		return err
	}

	fmt.Println("New ID was created:", todoObject.Id)

	return nil
}

func GetTodoAll() ([]structs.Todo_struct, error) {
	var allTodos []structs.Todo_struct

	ctx := context.Background()

	// Check database connection
	err := Db.PingContext(ctx)
	if err != nil {
		return allTodos, err
	}

	// The sql query to be executed
	tsql := "SELECT Id, Userid, Title, Category, State FROM [dbo].[todo];"

	// Execute query
	rows, err := Db.QueryContext(ctx, tsql)
	if err != nil {
		return allTodos, err
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
			return allTodos, err
		}

		todoItem := structs.Todo_struct{
			Id:       id,
			Userid:   userid,
			Title:    title,
			Category: category,
			State:    state,
		}

		allTodos = append(allTodos, todoItem)

		fmt.Printf("Id:%d Userid:%s Title: %s, Category: %s, State: %s\n", id, userid, title, category, state)
	}

	return allTodos, nil
}

func GetTodoObject(todoObject structs.Todo_struct) error {
	ctx := context.Background()

	// Check database connection
	err := Db.PingContext(ctx)
	if err != nil {
		return err
	}

	// The sql query to be executed
	tsql := "SELECT Id, Userid, Title, Category, State FROM [dbo].[todo] WHERE Userid = @Userid;"

	// Execute query
	rows, err := Db.QueryContext(ctx, tsql, sql.Named("UserId", todoObject.Userid))
	if err != nil {
		return err
	}

	// Wait until function ends before
	defer rows.Close()

	fmt.Println("Print one todo object:")

	// Iterate through the result set.
	for rows.Next() {

		// Get values from row.
		err := rows.Scan(&todoObject.Id,
			&todoObject.Userid,
			&todoObject.Title,
			&todoObject.Category,
			&todoObject.State)

		if err != nil {
			return err
		}

		fmt.Printf("Id:%d Userid:%s Title: %s, Category: %s, State: %s\n",
			todoObject.Id, todoObject.Userid, todoObject.Title, todoObject.Category, todoObject.State)
	}

	return nil
}
