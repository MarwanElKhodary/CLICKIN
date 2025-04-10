package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ** This structure found from: https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go
// func setupTest() func() {

// 	return func() {

//		}
//	}

func TestRootRoute(t *testing.T) {
	//TODO: Move to test setup function
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL database!")

	repo := NewRepository(config.DB)

	service := NewService(repo)

	handler := NewHandler(service)

	router := gin.Default()
	handler.SetupRoutes(router)
	//END OF TEST SETUP
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
