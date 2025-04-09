package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
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
	//TODO: Move this somewhere else
	router.LoadHTMLFiles("../frontend/index.html")
	router.Static("/css", "../frontend/css/")
	router.Static("/js", "../frontend/js/")
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{}, //Used to add headers
		)
	})
	router.Use(cors.Default())
	handler.SetupRoutes(router)
	router.Run("localhost:8080")
}
