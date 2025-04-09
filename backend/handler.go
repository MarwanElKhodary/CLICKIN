package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Handler contains the HTTP handlers and their dependencies
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.SetTrustedProxies(nil)
	router.Use(cors.Default())

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

	router.NoRoute(func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{}, //Used to add headers
		)
	})
	// ? Maybe change the name to /get-count and /update-count
	router.GET("/count", h.getCountHandler)
	router.POST("/count", h.incrementCountHandler)
}

// ? How come this function doesn't start with a capital letter but others do?
func (h *Handler) incrementCountHandler(c *gin.Context) {
	_, err := h.service.IncrementRandomCount()

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	count, err := h.service.GetTotalCount()

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d", count))
}

func (h *Handler) getCountHandler(c *gin.Context) {
	count, err := h.service.GetTotalCount()

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d", count))
}
