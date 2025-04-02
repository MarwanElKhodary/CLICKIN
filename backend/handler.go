package main

import (
	"fmt"
	"net/http"

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
	// ? Maybe change the name to /get-count and /update-count
	router.GET("/count", h.getCountHandler)
	router.POST("/count", h.incrementCountHandler)
}

// ? How come this function doesn't start with a capital letter but others do?
// ! A 204 OPTIONS call happens everytime there's a post request
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
