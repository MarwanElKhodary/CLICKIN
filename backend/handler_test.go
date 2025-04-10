package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ? What is this t * testing.T format?
// ? Read more about defer

var router *gin.Engine // So that you can use router in both setup and test methods

// ** This structure found from: https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go
func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")

	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo := NewRepository(config.DB) // TODO: Substitute with mock database

	service := NewService(repo)

	handler := NewHandler(service)

	router = gin.Default()
	handler.SetupRoutes(router)

	return func(t *testing.T) {
		t.Log("teardown test case")
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
		{"getCount", "GET", "/count", nil, 200},
		{"postCount", "POST", "/count", nil, 200}, // ! This actually increments the count in the database
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	w := httptest.NewRecorder()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, tc.url, tc.body)
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expected, w.Code)
		})
	}
}
