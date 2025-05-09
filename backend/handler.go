// Package main implements a simple click-counter application with a web interface.
// It provides a REST API to get and increment a counter, and a real-time
// WebSocket connection to update all clients when the counter changes.
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//TODO:

// Handler contains the HTTP handlers and their dependencies.
// It manages HTTP routes and translates between HTTP requests/responses
// and the application's service layers and their dependencies
type Handler struct {
	service *Service
}

// NewHandler creates a new handler with the given service.
// It initializes the handler with the service dependency needed to process requests.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// SetupRoutes configures all the routes for the application.
// It sets up static file serving, HTML rendering, API endpoints,
// and the WebSocket connection handler.
func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.SetTrustedProxies(nil)
	router.Use(cors.Default())

	router.LoadHTMLFiles("../frontend/index.html")
	router.Static("/assets", "../frontend/assets/")
	router.Static("/js", "../frontend/js/")

	// Versioning is for sanity and bypassing cache
	router.StaticFile("favicon_v3.ico", "../frontend/favicon_v3.ico")
	router.StaticFile("icon_v2.svg", "../frontend/icon_v2.svg")
	router.StaticFile("icon_v3.png", "../frontend/icon_v3.png")

	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{}, //Used to add headers
		)
	})

	router.GET("/robots.txt", func(c *gin.Context) {
		c.File("../frontend/robots.txt")
	})

	router.GET("/ws", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request)
	})

	router.NoRoute(func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{}, //Used to add headers
		)
	})

	// Added to mitigate css issues if user navigates to the url below
	// Current known issue is when the user goes to /assets/ it shows a blank screen
	router.GET("/frontend/*any", func(c *gin.Context) {
		// ? Not sure if this is even an issue
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	router.GET("/count", h.getCountHandler)
	router.POST("/count", h.incrementCountHandler)
}

// incrementCountHandler handles POST requests to increment the counter.
// It increments a random counter slot, broadcasts the count using WebSockets using the lastInsertId.
// Returns HTTP 400 Bad Request if there is an error during increment.
func (h *Handler) incrementCountHandler(c *gin.Context) {
	lastInsertId, err := h.service.IncrementRandomCount()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	BroadcastCount(int(lastInsertId)) // ? Not sure if this will always work
	c.Status(http.StatusOK)
}

// getCountHandler handles GET requests to retrieve the current counter value.
// It returns the total count across all counter slots.
// Returns HTTP 404 Not Found if there is an error retrieving the count.
func (h *Handler) getCountHandler(c *gin.Context) {
	count, err := h.service.GetTotalCount()

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d", count))
}
