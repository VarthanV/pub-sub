package server

import "github.com/gin-gonic/gin"

func InitRoutes(r *gin.Engine, ctrl *Controller) {

	wsGroup := r.Group("/ws")
	exchangeGroup := r.Group("/exchange")

	wsGroup.GET("", ctrl.HandleSubscription)

	exchangeGroup.POST("", ctrl.CreateExchange)
}
