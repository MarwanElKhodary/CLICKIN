// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ? What is this t * testing.T format?
// ? Read more about global variables in golang

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

// TestRoutes verifies that all defined routes return the expected status codes.
// It tests the root route, the get count endpoint, and the post count endpoint.
func TestRoutes(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		url      string
		body     io.Reader
		expected int
	}{
		{"root", "GET", "/", nil, http.StatusOK},
		{"getCount", "GET", "/count", nil, http.StatusOK},
		{"postCount", "POST", "/count", nil, http.StatusOK},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	mock.MatchExpectationsInOrder(false)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO count_table (slot, count) VALUES (?, ?)")).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := mock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(count) as count FROM count_table")).
		WillReturnRows(rows)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, tc.body)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

// TestSameSlots tests concurrent incrementing of the same counter slot.
// It verifies that 100 concurrent increments to the same slot result in
// the expected total count.
func TestSameSlots(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	numRequests := 100
	sameSlot := 69

	for i := range numRequests {
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO count_table (slot, count) VALUES (?, ?)")).
			WithArgs(sameSlot, 1).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}
	t.Run("sameSlotCollision", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(numRequests)

		for range numRequests {
			go func() {
				defer wg.Done()

				_, err := repo.IncrementCount(sameSlot, 1)
				if err != nil {
					t.Errorf("Failed to increment: %v", err)
				}
			}()
		}

		wg.Wait()

		rows := mock.NewRows([]string{"count"}).AddRow(numRequests)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(count) as count FROM count_table")).
			WillReturnRows(rows)

		totalCount, err := repo.GetTotalCount()
		if err != nil {
			t.Fatalf("Failed to get count: %v", err)
		}

		assert.Equal(t, numRequests, totalCount)
	})
}
