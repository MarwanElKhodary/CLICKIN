package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ? What is this t * testing.T format?
// ? Read more about global variables in golang
// ? Read more about `if err := mock.ExpectationsWereMet(); err != nil {` notation

var router *gin.Engine
var mock sqlmock.Sqlmock

// ** This structure found from: https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go
func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")

	db, mockSQL, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock = mockSQL

	repo := NewRepository(db)

	service := NewService(repo)

	handler := NewHandler(service)

	router = gin.Default()
	handler.SetupRoutes(router)

	return func(t *testing.T) {
		t.Log("teardown test case")
		db.Close()
	}
}

// TODO: Consider refactoring this into separate methods instead
func TestRoutes(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		url      string
		body     io.Reader
		expected int
	}{
		{"root", "GET", "/", nil, 200},
		{"getCount", "GET", "/count", nil, 200},
		{"postCount", "POST", "/count", nil, 200},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	mock.MatchExpectationsInOrder(false)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO count_table (slot, count) VALUES (?, ?)")).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Then: Expect two SELECTs (one for POST /count, one for GET /count)
	rows := mock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(count) as count FROM count_table")).
		WillReturnRows(rows)
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
