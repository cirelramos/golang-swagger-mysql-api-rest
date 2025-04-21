package routes

import (
	"database/sql"
	"github.com/gorilla/mux"
	"golang-api-rest-swagger/controllers" // Import the models package
	"net/http"
)

// SetupRoutes defines the API routes and associates them with the appropriate handler functions.
func SetupRoutes(r *mux.Router, db *sql.DB) { // Add db as parameter
	r.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetBooks(w, r, db)
	}).Methods("GET")

	r.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetBook(w, r, db)
	}).Methods("GET")

	r.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateBook(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateBook(w, r, db)
	}).Methods("PUT")

	r.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteBook(w, r, db)
	}).Methods("DELETE")
}
