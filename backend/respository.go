package main

import (
	"database/sql"
	"fmt"
)

// Repository handles the database operations
type Repository struct {
	db *sql.DB
}

// TODO: Investigate if this is actually necessary
// I guess it makes sense in the future, but not so much right now
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetTotalCount() (int, error) {
	var cnt int
	row := r.db.QueryRow("SELECT SUM(count) as count FROM count_table")
	err := row.Scan(&cnt)

	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("GetTotalCount: couldn't get count")
	}

	return cnt, nil
}

func (r *Repository) IncrementCount(slot int, count int) (int64, error) {
	result, err := r.db.Exec("INSERT INTO count_table (slot, count) VALUES (?, ?)", slot, count)

	if err != nil {
		return 0, fmt.Errorf("IncrementCount: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("IncrementCount: %v", err)
	}

	return id, nil
}
