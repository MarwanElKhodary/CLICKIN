package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ? What is this t * testing.T format?
// ? Read more about defer
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

func TestRootRoute(t *testing.T) {
	cases := []struct {
		name     string
		method   string
		url      string
		body     io.Reader
		expected int
	}{
		{"root", "GET", "/", nil, 200},
		{"getCount", "GET", "/count", nil, 200},   // ! Currently failing
		{"postCount", "POST", "/count", nil, 200}, // ! Currently failing
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(tc.method, tc.url, tc.body)
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expected, w.Code)
			if err := mock.ExpectationsWereMet(); err != nil { // ! Currently doesn't do anything
				t.Errorf("not all expectations were met: %v", err)
			}
		})
	}
}
