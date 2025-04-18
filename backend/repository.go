// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"database/sql"
	"fmt"
)

// Repository handles database operations for the counter application.
// It provides methods to retrieve and update counter data in the database.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository with the given database connection.
// The repository is responsible for all database interactions.
//
// TODO: Look into using an ORM for security
// TODO: Investigate using the context package and passing it here for mock tests
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetTotalCount retrieves the sum of all counts from the database.
// It returns the total count and any error encountered.
func (r *Repository) GetTotalCount() (int, error) {
	var cnt int
	row := r.db.QueryRow("SELECT SUM(count) as count FROM count_table")
	err := row.Scan(&cnt)

	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("repository.GetTotalCount: %v", err)
	}

	return cnt, nil
}

// IncrementCount increments the counter for a specific slot.
// It inserts a new row with the given slot and count values,
// and returns the last insert ID and any error encountered.
func (r *Repository) IncrementCount(slot int, count int) (int64, error) {
	result, err := r.db.Exec("INSERT INTO count_table (slot, count) VALUES (?, ?)", slot, count)

	if err != nil {
		return 0, fmt.Errorf("repository.IncrementCount: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("repository.IncrementCount: %v", err)
	}

	return id, nil
}
