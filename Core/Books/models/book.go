package models

// Book struct to hold book details.
type Book struct {
	ID     int    `json:"id" db:"id"`
	Title  string `json:"title" db:"title"`
	Author string `json:"author" db:"author"`
	Year   int    `json:"year" db:"year"`
}
