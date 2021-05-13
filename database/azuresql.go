package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"projectGroup23/structs"
)

var Db *sql.DB

var Err error

func UpdateTodoObject(sqlId int, description string) error {
	ctx := context.Background()

	if Db == nil {
		Err = errors.New("create todo object: db is null")
	}

	// Check db connection
	Err = Db.PingContext(ctx)
	if Err != nil {
		return Err
	}

	tsql := "UPDATE [dbo].[todo] SET Description = @Description WHERE Id = @Id;"

	// Execute non-query with named parameters
	_, err := Db.ExecContext(ctx, tsql,
		sql.Named("Id", sqlId),
		sql.Named("Description", description))

	if err != nil {
		return err
	}

	return nil
}

func UpdateTodoObjectStatus(sqlId int, status string) error {
	ctx := context.Background()

	if Db == nil {
		Err = errors.New("create todo object: db is null")
	}

	// Check db connection
	Err = Db.PingContext(ctx)
	if Err != nil {
		return Err
	}

	tsql := "UPDATE [dbo].[todo] SET State = @State WHERE Id = @Id;"

	// Execute non-query with named parameters
	_, err := Db.ExecContext(ctx, tsql,
		sql.Named("Id", sqlId),
		sql.Named("State", status))

	if err != nil {
		return err
	}

	return nil
}

func DeleteTodoObject(sqlId int) error {
	ctx := context.Background()

	if Db == nil {
		Err = errors.New("create todo object: db is null")
	}

	// Check db connection
	Err = Db.PingContext(ctx)
	if Err != nil {
		return Err
	}

	tsql := "DELETE FROM [dbo].[todo] WHERE Id = @Id;"

	// Execute non-query with named parameters
	res, err := Db.ExecContext(ctx, tsql,
		sql.Named("Id", sqlId))

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

	tsql := `INSERT INTO [dbo].[todo] (Userid, Description, State) VALUES (@Userid, @Description, @State);
			select isNull(SCOPE_IDENTITY(), -1);`

	stmt, err := Db.Prepare(tsql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx,
		sql.Named("Userid", todoObject.Userid),
		sql.Named("Description", todoObject.Description),
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
	tsql := "SELECT Id, Userid, Description, State FROM [dbo].[todo];"

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
		var userid, description, state string
		var id int

		// Get values from row.
		err := rows.Scan(&id, &userid, &description, &state)
		if err != nil {
			return allTodos, err
		}

		todoItem := structs.Todo_struct{
			Id:          id,
			Userid:      userid,
			Description: description,
			State:       state,
		}

		allTodos = append(allTodos, todoItem)

		fmt.Printf("Id:%d Userid:%s Description: %s, State: %s\n", id, userid, description, state)
	}

	return allTodos, nil
}

func GetTodoObject(userid string) ([]structs.Todo_struct, error) {
	var allTodos []structs.Todo_struct

	ctx := context.Background()

	// Check database connection
	err := Db.PingContext(ctx)
	if err != nil {
		return allTodos, err
	}

	// The sql query to be executed
	tsql := "SELECT Id, Userid, Description, State FROM [dbo].[todo] WHERE Userid = @Userid;"

	// Execute query
	rows, err := Db.QueryContext(ctx, tsql, sql.Named("Userid", userid))
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

		todoItem := structs.Todo_struct{
			Id:          id,
			Userid:      userid,
			Description: description,
			State:       state,
		}

		allTodos = append(allTodos, todoItem)
	}

	return allTodos, nil
}
