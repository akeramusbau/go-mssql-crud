package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

type User struct {
	ID       int
	Username string
	Email    string
}

func main() {
	db, err := connectToSQLServer()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1:
		{
			// Create a user
			pk, _ := insertData(db, "akera musbau", "akeramusbau@gmail.com")
			fmt.Printf("ID : %d", pk)
		}
	case 2:
		{
			// Update a user
			updatedRows, err := updateData(db, 14, "akera_updated", "akeramusbau@gmail.com")
			if err != nil {
				panic(err)
			}

			fmt.Printf("Number of updated rows: %d", updatedRows)
		}
	case 3:
		{
			// Delete a user
			deletedRows, err := deleteData(db, 14)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Number deleted users: %d", deletedRows)
		}
	default:
		{
			// Read users
			users, err := queryData(db)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Users:")
			for _, user := range users {
				fmt.Printf("ID: %d, Username: %s, Email: %s\n", user.ID, user.Username, user.Email)
			}
		}
	}

}
func connectToSQLServer() (*sql.DB, error) {
	server := "servername"
	port := 1433
	user := "username"
	password := "password"
	database := "database"

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	// Ping the SQL Server to ensure connectivity
	err = db.PingContext(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to SQL Server!")
	return db, nil
}

func queryData(db *sql.DB) ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
func updateData(db *sql.DB, id int, newUsername, newEmail string) (int64, error) {
	stmt, err := db.Prepare("UPDATE users SET username = @p1, email = @p2 WHERE ID = @P3")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUsername, newEmail, id)
	if err != nil {
		return 0, nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
func deleteData(db *sql.DB, id int) (int64, error) {
	stmt, err := db.Prepare("DELETE FROM Users WHERE ID = @p1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
func insertData(db *sql.DB, username, email string) (int, error) {
	stmt, err := db.Prepare("INSERT INTO users(USERNAME,EMAIL) OUTPUT INSERTED.id VALUES (@p1,@P2)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var lastInsertedId int
	err = stmt.QueryRow(username, email).Scan(&lastInsertedId)
	if err != nil {
		return 0, nil
	}

	return lastInsertedId, nil
}
