package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang-api-rest-swagger/models" // Import the models package
	"net/http"
	"strconv"
)

// GetBooks handles the retrieval of all books from the database.
// @Summary Get all books
// @Description Retrieve a list of all books from the database
// @Tags books
// @Produce json
// @Success 200 {array} models.Book
// @Router /books [get]
func GetBooks(w http.ResponseWriter, r *http.Request, db *sql.DB) { // Add db as parameter
	w.Header().Set("Content-Type", "application/json")

	// Query the database.
	rows, err := db.Query("SELECT id, title, author, YEAR FROM books")
	if err != nil {
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the results.
	books := []models.Book{} // Use models.Book

	// Iterate over the rows.
	for rows.Next() {
		var book models.Book // Use models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year); err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan row: %v", err), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error during row iteration: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}

// GetBook handles the retrieval of a single book by ID from the database.
// @Summary Get a book by ID
// @Description Retrieve a single book by its ID from the database
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [get]
func GetBook(w http.ResponseWriter, r *http.Request, db *sql.DB) { // Add db as parameter
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	// Query the database for the book with the given ID.
	row := db.QueryRow("SELECT id, title, author, YEAR FROM books WHERE id = ?", id)
	var book models.Book // Use models.Book
	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Database query failed: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// CreateBook handles the creation of a new book in the database.
// @Summary Create a new book
// @Description Add a new book to the database
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book object to be added"
// @Success 201 {object} models.Book
// @Failure 400 {string} string "Invalid request body"
// @Router /books [post]
func CreateBook(w http.ResponseWriter, r *http.Request, db *sql.DB) { // Add db as parameter
	w.Header().Set("Content-Type", "application/json")
	var book models.Book // Use models.Book
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
	book.ID = int(insertID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// UpdateBook handles the updating of an existing book in the database.
// @Summary Update an existing book
// @Description Update the details of an existing book in the database
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Book true "Updated book object"
// @Success 200 {object} models.Book
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [put]
func UpdateBook(w http.ResponseWriter, r *http.Request, db *sql.DB) { // Add db as parameter
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updatedBook models.Book // Use models.Book
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

// DeleteBook handles the deletion of a book from the database.
// @Summary Delete a book
// @Description Delete a book from the database
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {string} string "Book deleted successfully"
// @Failure 404 {string} string "Book not found"
// @Router /books/{id} [delete]
func DeleteBook(w http.ResponseWriter, r *http.Request, db *sql.DB) { // Add db as parameter.
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
