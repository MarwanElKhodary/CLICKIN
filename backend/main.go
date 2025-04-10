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
	handler.SetupRoutes(router) // ! The router is dependent on the handler at this current stage
	//TODO: Consider separating the router logic to make testing easier
	router.Run("localhost:8080")
}
