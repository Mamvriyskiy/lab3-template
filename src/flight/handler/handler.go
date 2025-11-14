package handler

import (
	"net/http"

	services "github.com/Mamvriyskiy/lab3-template/src/flight/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	router.GET("/manage/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	flight := router.Group("")

	flight.GET("/flight", h.GetInfoAboutFlight)
	flight.GET("/flight/:flightNumber", h.GetInfoAboutFlightByFlightNumber)

	return router
}
