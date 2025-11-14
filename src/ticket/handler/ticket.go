package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateStatusTicket(c *gin.Context) {
    ticketUID := c.Param("ticketUid")
	if ticketUID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ticketUid is required"})
        return
    }

    err := h.services.UpdateStatusTicket(ticketUID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusOK)
}

func (h *Handler) GetInfoAboutTiket(c *gin.Context) {
	ticketUID := c.Param("ticketUid")
	if ticketUID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ticketUid is required"})
        return
    }

	resp, err := h.services.GetInfoAboutTiket(ticketUID)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetInfoAboutTikets(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
    if username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
        return
    }

	resp, err := h.services.GetInfoAboutTikets(username)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

type BuyTicket struct {
    FlightNumber    string `json:"flightNumber"`
    Price           int    `json:"price"`
    PaidFromBalance bool   `json:"paidFromBalance"`
}

func (h *Handler) CreateTicket(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
    if username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
        return
    }

    var buyTicket BuyTicket
	if err := c.BindJSON(&buyTicket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.services.CreateTicket(username, buyTicket.FlightNumber, buyTicket.Price)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}
