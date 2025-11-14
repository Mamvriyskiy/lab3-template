package handler

import (
	"net/http"
	"time"

	cb "github.com/Mamvriyskiy/lab3-template/src/gateway/circuitBreaker"
	services "github.com/Mamvriyskiy/lab3-template/src/gateway/services"
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

	// Создаем CB
	readCB := &cb.CircuitBreaker{
		FailureThreshold: 3,
		RetryTimeout:     5 * time.Second,
	}

	router.Use(cb.NewCBMiddleware(readCB, cb.FallbackHandler))

	router.GET("/manage/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	flight := router.Group("api/v1/")

	// READ endpoints (идут через circuit breaker)
	flight.GET("/flights", h.GetInfoAboutFlight)
	flight.GET("/me", h.GetInfoAboutUser)
	flight.GET("/tickets", h.GetInfoAboutAllUserTickets)
	flight.GET("/tickets/:ticketUid", h.GetInfoAboutUserTicket)
	flight.GET("/privilege", h.GetInfoAboutUserPrivilege)

	// WRITE (не трогаем)
	flight.POST("/tickets", h.BuyTicketUser)
	flight.DELETE("/tickets/:ticketUid", h.DeleteTicketUSer)

	return router
}
