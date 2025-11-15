package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


func (h *Handler) GetInfoAboutUserPrivilege(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
    if username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
        return
    }

    resp, err := h.services.GetInfoAboutUserPrivilege(username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateBonusBonus(c *gin.Context) {
    username := c.GetHeader("X-User-Name")
    if username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
        return
    }

    priceStr := c.Param("price")

    price, err := strconv.Atoi(priceStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ticketUid := c.Param("ticketUID")
    if ticketUid == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ticketUid is required"})
        return
    }

    updateBonus, _ := h.services.UpdateBonusBonus(username, ticketUid, price)

    c.JSON(http.StatusOK, gin.H{
        "updated_balance": updateBonus,
    })
}


func (h *Handler) UpdateBonus(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
    if username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
        return
    }

    ticketUid := c.Param("ticketUID")
    if ticketUid == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ticketUid is required"})
        return
    }

    priceStr := c.Param("price")

    price, err := strconv.Atoi(priceStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if price == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    info, err := h.services.UpdateBonus(username, ticketUid, price)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, info)
}

func (h *Handler) UpdateBonusDelete(c *gin.Context) {
    username := c.GetHeader("X-User-Name")
    if username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
        return
    }
    
    priceStr := c.Param("price")

    price, err := strconv.Atoi(priceStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
  
    h.services.UpdateBonusDelete(username, price)
    c.Status(http.StatusOK)
}
