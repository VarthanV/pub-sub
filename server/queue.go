package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type createQueueRequest struct {
	Name    string `json:"name" binding:"required"`
	Durable bool   `json:"durable"`
}

func (c *Controller) CreateQueue(ctx *gin.Context) {
	request := createQueueRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		logrus.Error("unable to bind request ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = c.broker.CreateQueue(ctx, request.Name, request.Durable)
	if err != nil {
		logrus.Error("error in creating queue ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}
