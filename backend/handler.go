package main

import (
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
	router.GET("/count", h.getCountHandler)
	router.POST("/count", h.incrementCountHandler)
}

// ? How come this function doesn't start with a capital letter but others do?
func (h *Handler) incrementCountHandler(c *gin.Context) {
	id, err := h.service.IncrementRandomCount()

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, id)
	}
}

func (h *Handler) getCountHandler(c *gin.Context) {
	count, err := h.service.GetTotalCount()

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, count)
	}
}
