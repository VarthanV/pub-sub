package server

import "github.com/gin-gonic/gin"

func InitRoutes(r *gin.Engine, ctrl *Controller) {

	wsGroup := r.Group("/ws")
	exchangeGroup := r.Group("/exchanges")
	queueGroup := r.Group("/queues")
	bindingGroup := r.Group("/bindings")
	publishGroup := r.Group("/publish")

	wsGroup.GET("", ctrl.HandleSubscription)

	exchangeGroup.POST("", ctrl.CreateExchange)

	queueGroup.POST("", ctrl.CreateQueue)

	bindingGroup.POST("", ctrl.CreateBinding)

	publishGroup.POST("", ctrl.Publish)
}
