// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL database!")

	repo := NewRepository(config.DB)

	service := NewService(repo)

	handler := NewHandler(service)

	router := gin.Default()

	// ! The router is dependent on the handler at this current stage
	//TODO: Consider separating the router logic to make testing easier
	handler.SetupRoutes(router)

	// ! Super inconcistent when using on a different device on the same wifi
	router.Run("0.0.0.0:8080")
}
