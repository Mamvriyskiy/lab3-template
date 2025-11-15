package handler

import (
	"net/http"

	services "github.com/Mamvriyskiy/lab3-template/src/bonus/services"
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

	bonus := router.Group("")

	bonus.GET("/privilege", h.GetInfoAboutUserPrivilege)
	bonus.PATCH("/bonus/:ticketUID/:price", h.UpdateBonus)
	bonus.PATCH("/bonusUpdate/:ticketUID/:price", h.UpdateBonusBonus)
	bonus.DELETE("/bonusUpdateDelete/:price", h.UpdateBonusDelete)

	return router
}
