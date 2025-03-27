package main

import "math/rand/v2"

// Service contains the business logic of the application
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// ? Again, I don't really understand the point of this
func (s *Service) GetTotalCount() (int, error) {
	return s.repo.GetTotalCount()
}

func (s *Service) IncrementRandomCount() (int64, error) {
	return s.repo.IncrementCount(rand.IntN(100), 1)
}
