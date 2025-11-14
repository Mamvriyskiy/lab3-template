package handler

import (
	"net/http"

	services "github.com/Mamvriyskiy/lab3-template/src/ticket/services"
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

	ticket := router.Group("")

	ticket.GET("/ticket/:ticketUid", h.GetInfoAboutTiket)
	ticket.GET("/tickets", h.GetInfoAboutTikets)
	ticket.PATCH("/ticket/:ticketUid", h.UpdateStatusTicket)
	ticket.POST("/ticket", h.CreateTicket)

	return router
}
