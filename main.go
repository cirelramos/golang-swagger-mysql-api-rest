package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // Import the godotenv package
	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag/example/celler/docs" // Import the generated docs
	"os"
)

// Book struct to hold book details.  The `json` tags are for JSON serialization,
// the `db` tags are for gorm (if you were to use an ORM).
type Book struct {
	ID     int    `json:"id" db:"id"`
	Title  string `json:"title" db:"title"`
	Author string `json:"author" db:"author"`
	Year   int    `json:"year" db:"year"`
}

// Global variable to hold the database connection.
var db *sql.DB

// initDB initializes the database connection.  It should be called only once
// during the application's startup.
func initDB() {
	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve database credentials from environment variables.
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")

	// Check if the environment variables are set.
	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatalf("Database credentials not set in .env file")
	}

	// Construct the connection string.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Connect to the database.
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set maximum number of connections.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	// Check if the connection is working.
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to MySQL database!")

	// Create the books table if it doesn't exist.
	// IMPORTANT:  Use a proper migration tool (like https://github.com/golang-migrate/migrate)
	//             for production databases.  This is just for demonstration.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			YEAR INT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

// Handlers

// Get all books from the database
// @Summary Get all books
// @Description Retrieve a list of all books from the database
// @Tags books
// @Produce json
// @Success 200 {array} Book
// @Router /books [get]
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query the database.  Use a SELECT statement.
	rows, err := db.Query("SELECT id, title, author, YEAR FROM books")
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Ensure rows are closed after we're done with them.

	// Create a slice to hold the results.
	books := []Book{}

	// Iterate over the rows.
	for rows.Next() {
		var book Book
		// Use rows.Scan to populate the Book struct.  The order of the arguments
		// MUST match the order of the columns in the SELECT query.
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year); err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan row: %v", err), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	// Check for errors during row iteration.
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error during row iteration: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode the results as JSON.
	json.NewEncoder(w).Encode(books)
}

// Get a single book by ID from the database
// @Summary Get a book by ID
// @Description Retrieve a single book by its ID from the database
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} Book
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [get]
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	// Query the database for the book with the given ID.
	row := db.QueryRow("SELECT id, title, author, YEAR FROM books WHERE id = ?", id)
	var book Book
	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode the book as JSON.
	json.NewEncoder(w).Encode(book)
}

// Create a new book in the database
// @Summary Create a new book
// @Description Add a new book to the database
// @Tags books
// @Accept json
// @Produce json
// @Param book body Book true "Book object to be added"
// @Success 201 {object} Book
// @Failure 400 {string} string "Invalid request body"
// @Router /books [post]
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	if book.Title == "" || book.Author == "" || book.Year == 0 {
		http.Error(w, "Invalid request body: Title, Author, and Year are required", http.StatusBadRequest)
		return
	}

	// Insert the new book into the database.
	result, err := db.Exec("INSERT INTO books (title, author, year) VALUES (?, ?, ?)", book.Title, book.Author, book.Year)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database insert failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly inserted book.
	insertID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get last insert ID: %v", err), http.StatusInternalServerError)
		return
	}
	book.ID = int(insertID) // Convert int64 to int.

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// Update an existing book in the database
// @Summary Update an existing book
// @Description Update the details of an existing book in the database
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body Book true "Updated book object"
// @Success 200 {object} Book
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [put]
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	if updatedBook.Title == "" || updatedBook.Author == "" || updatedBook.Year == 0 {
		http.Error(w, "Invalid request body: Title, Author, and Year are required", http.StatusBadRequest)
		return
	}

	// Update the book in the database.
	result, err := db.Exec("UPDATE books SET title = ?, author = ?, year = ? WHERE id = ?", updatedBook.Title, updatedBook.Author, updatedBook.Year, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database update failed: %v", err), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get number of updated rows: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	updatedBook.ID = id
	json.NewEncoder(w).Encode(updatedBook)
}

// Delete a book from the database
// @Summary Delete a book
// @Description Delete a book from the database
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {string} string "Book deleted successfully"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [delete]
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	// Delete the book from the database.
	result, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database delete failed: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get number of deleted rows: %v", err), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted successfully"})
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server for a book management API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Initialize the database connection.
	initDB()
	defer db.Close() // Ensure the connection is closed when the program exits.

	// Create a new router
	r := mux.NewRouter()

	// Define API endpoints.
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// Swagger documentation endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
