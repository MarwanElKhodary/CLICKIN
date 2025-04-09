package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// TODO: Add file server here
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
	handler.SetupRoutes(router)
	router.Run("localhost:8080")
}
