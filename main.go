package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
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
	pk, _ := insertData(db, user.Username, user.Email, user.Password)
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

	// Apply middleware globally
	r.Use(loggingMiddleware)
	r.Use(errorHandlingMiddleware)

	// Define the endpoints
	r.HandleFunc("/login", login).Methods("POST")
	r.Handle("/users", authenticate(http.HandlerFunc(getUsers))).Methods("GET")
	r.Handle("/users/{id}", authenticate(http.HandlerFunc(getUser))).Methods("GET")
	r.Handle("/users", authenticate(http.HandlerFunc(createUser))).Methods("POST")
	r.Handle("/users/{id}", authenticate(http.HandlerFunc(updateUser))).Methods("PUT")
	r.Handle("/users/{id}", authenticate(http.HandlerFunc(deleteUser))).Methods("DELETE")

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
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
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
func insertData(db *sql.DB, username, email string, password string) (int, error) {
	stmt, err := db.Prepare("INSERT INTO users(USERNAME,EMAIL,password) OUTPUT INSERTED.id VALUES (@p1,@P2,@p3)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var lastInsertedId int
	err = stmt.QueryRow(username, email, password).Scan(&lastInsertedId)
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
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func userLogin(db *sql.DB, username string, email string) ([]User, error) {
	query := "SELECT * FROM users where username = @p1 and password = @p2"
	rows, err := db.QueryContext(context.Background(), query, username, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// authentication
var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}
func login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := connectToSQLServer()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var users []User
	users, err = userLogin(db, creds.Username, creds.Password)
	if err != nil {
		log.Fatal(err)
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := generateToken(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(5 * time.Minute),
	})
}
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the method and the requested URL
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log how long it took
		log.Printf("Completed in %v", time.Since(start))
	})
}
func errorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and send a user-friendly message
				log.Printf("Error occurred: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
