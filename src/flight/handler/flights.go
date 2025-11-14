package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)


func (h *Handler) GetInfoAboutFlight(c *gin.Context) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err := strconv.Atoi(pageStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }

    size, err := strconv.Atoi(sizeStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }


	flightsList, err := h.services.GetInfoAboutFlight(page, size)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, flightsList)
}

func (h *Handler) GetInfoAboutFlightByFlightNumber(c *gin.Context) {
    flightNumber := c.Param("flightNumber")
    if flightNumber == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Empty flightNumer"})
    }

    flight, err := h.services.GetInfoAboutFlightByFlightNumber(flightNumber)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, flight)
} 
