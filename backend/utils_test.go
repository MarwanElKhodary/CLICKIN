// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

type TestSuite struct {
	Db      *sql.DB
	SqlMock sqlmock.Sqlmock
	Router  *gin.Engine
	Repo    *Repository
	Serv    *Service
	Handler *Handler
}

var testSuite *TestSuite

// TestMain is called once before any tests are run
func TestMain(m *testing.M) {
	var teardown func()
	testSuite, teardown = setupTestSuite()

	code := m.Run()

	teardown()

	os.Exit(code)
}

// setupTestSuite initializes test dependencies and returns a TestSuite and teardown function.
// It creates a mock database, repository, service, and handler, and sets up routes.
func setupTestSuite() (*TestSuite, func()) {
	db, mockSql, err := sqlmock.New()
	if err != nil {
		panic("Failed to create a mock database: " + err.Error())
	}

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)
	router := gin.Default()
	handler.SetupRoutes(router)

	suite := &TestSuite{
		Db:      db,
		SqlMock: mockSql,
		Router:  router,
		Repo:    repo,
		Serv:    service,
		Handler: handler,
	}

	return suite, func() {
		db.Close()
	}
}
