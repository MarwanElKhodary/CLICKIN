// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// Global variables for testing
var router *gin.Engine
var mock sqlmock.Sqlmock
var repo *Repository

// setupTestCase initializes test dependencies and returns a teardown function.
// It creates a mock database, repository, service, and handler, and sets up routes.
// The returned function should be deferred to clean up resources after the test.
//
// This structure is based on: https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go
func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")

	db, mockSQL, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock = mockSQL

	repo = NewRepository(db)

	service := NewService(repo)

	handler := NewHandler(service)

	router = gin.Default()
	handler.SetupRoutes(router)

	return func(t *testing.T) {
		t.Log("teardown test case")
		db.Close()
	}
}
