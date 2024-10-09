package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

type User struct {
	ID       int
	Username string
	Email    string
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	// Read users
	db, err := connectToSQLServer()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	users, err := queryData(db)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(users)
}
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := connectToSQLServer()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	pk, _ := insertData(db, user.Username, user.Email)
	fmt.Printf("ID : %d", pk)
	json.NewEncoder(w).Encode("Successful")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	db, err := connectToSQLServer()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	i, err := strconv.Atoi(id)
	if err != nil {
		// ... handle error
		panic(err)
	}
	users, err := queryDataById(db, i)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(users)
}
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	db, err := connectToSQLServer()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	i, err := strconv.Atoi(id)
	if err != nil {
		// ... handle error
		panic(err)
	}
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	pk, _ := updateData(db, i, user.Username, user.Email)
	fmt.Printf("ID : %d", pk)
	json.NewEncoder(w).Encode("Record Updated Successful")
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	db, err := connectToSQLServer()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	i, err := strconv.Atoi(id)
	if err != nil {
		// ... handle error
		panic(err)
	}
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	pk, _ := deleteData(db, i)
	fmt.Printf("ID : %d", pk)
	json.NewEncoder(w).Encode("Record Deleted Successful")
}
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	fmt.Printf("Starting at port 8000\n")

	log.Fatal(http.ListenAndServe(":8000", r))

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
func queryDataById(db *sql.DB, id int) ([]User, error) {
	query := "SELECT * FROM users where id = @p1"
	rows, err := db.QueryContext(context.Background(), query, id)
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
