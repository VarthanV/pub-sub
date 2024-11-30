package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type publishRequest struct {
	ExchangeName string      `json:"exchange_name" binding:"required"`
	Key          string      `json:"key"`
	Payload      interface{} `json:"payload" binding:"required"`
}

func (c *Controller) Publish(ctx *gin.Context) {
	var (
		request = publishRequest{}
	)

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		logrus.Error("error in binding request ", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = c.broker.PublishMessage(ctx, request.ExchangeName, request.Key, request.Payload)
	if err != nil {
		logrus.Error("error in publishing message ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
