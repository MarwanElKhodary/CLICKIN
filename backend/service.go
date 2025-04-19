// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import "math/rand/v2"

// Service contains the business logic of the application.
// It coordinates operations between the repository and the handlers,
// implementing application-specific rules and behaviors.
type Service struct {
	repo *Repository
}

// NewService creates a new service with the given repository.
// The service encapsulates the core business logic of the application.
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetTotalCount returns the total count from all slots in the counter.
// It delegates to the repository to fetch the data.
func (s *Service) GetTotalCount() (int, error) {
	return s.repo.GetTotalCount()
}

// IncrementRandomCount increments the count for a random slot.
// It selects a random slot between 0 and 99, then increments it by 1.
// Returns the ID of the new record and any error encountered.
func (s *Service) IncrementRandomCount() (int64, error) {
	return s.repo.IncrementCount(rand.IntN(100), 1)
}
