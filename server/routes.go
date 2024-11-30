package server

import "github.com/gin-gonic/gin"

func InitRoutes(r *gin.Engine, ctrl *Controller) {

	wsGroup := r.Group("/ws")
	exchangeGroup := r.Group("/exchanges")
	queueGroup := r.Group("/queues")

	wsGroup.GET("", ctrl.HandleSubscription)

	exchangeGroup.POST("", ctrl.CreateExchange)

	queueGroup.POST("", ctrl.CreateQueue)
}
