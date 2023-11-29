package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

)

// Model struct
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {

	// Initialize router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getDb() *sql.DB {
	// Connect to MySQL
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/go")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MySQL")
	return db
}

// Get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	// Query the database
	log.Println("Get all users")
	db := getDb()
	rows, err := db.Query("SELECT * FROM users")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create a slice to hold the results
	var users []User

	// Iterate over the rows
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	defer db.Close()
	// Convert the slice to JSON
	json.NewEncoder(w).Encode(users)
}

// Get a single user
func getUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request parameters
	db := getDb()
	params := mux.Vars(r)
	id := params["id"]

	// Query the database
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	// Create a user object
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the user object to JSON
	json.NewEncoder(w).Encode(user)
}

// Create a user
func createUser(w http.ResponseWriter, r *http.Request) {
	db := getDb()
	// Create a new user object from the request body
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	// Insert the user into the database
	_, err := db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", user.Username, user.Email)
	if err != nil {
		log.Fatal(err)
	}

	// Return a success message
	json.NewEncoder(w).Encode("User created successfully")
}

// Update a user
func updateUser(w http.ResponseWriter, r *http.Request) {
	db := getDb()
	// Get the user ID from the request parameters
	params := mux.Vars(r)
	id := params["id"]

	// Create a new user object from the request body
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	// Update the user in the database
	_, err := db.Exec("UPDATE users SET username = ?, email = ? WHERE id = ?", user.Username, user.Email, id)
	if err != nil {
		log.Fatal(err)
	}

	// Return a success message
	json.NewEncoder(w).Encode("User updated successfully")
}

// Delete a user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	db := getDb()
	// Get the user ID from the request parameters
	params := mux.Vars(r)
	id := params["id"]

	// Delete the user from the database
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	// Return a success message
	json.NewEncoder(w).Encode("User deleted successfully")
}
