package main

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"golang-api-rest-swagger/database"
	_ "golang-api-rest-swagger/docs" // Import the generated docs
	"golang-api-rest-swagger/routes"
	"log"
	"net/http"
)

// main.go
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
	// Initialize database connection
	db, err := database.InitDB() // Changed to package call
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create a new router
	r := mux.NewRouter()

	// Define routes using the routes package
	routes.SetupRoutes(r, db) // Changed to package call

	// Swagger documentation endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start the server
	port := ":8080"
	log.Println("start in port " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
