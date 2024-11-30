package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type createBindingRequest struct {
	Key          string
	ExchangeName string `json:"exchange_name" binding:"required"`
	QueueName    string `json:"queue_name" binding:"required"`
}

func (c *Controller) CreateBinding(ctx *gin.Context) {
	var (
		request = createBindingRequest{}
	)

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		logrus.Error("error in binding request ", err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = c.broker.BindQueue(ctx, request.QueueName,
		request.ExchangeName, request.Key)
	if err != nil {
		logrus.Error("error in binding req ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
