package server

import (
	"net/http"

	"github.com/VarthanV/pub-sub/exchange"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CreateExchangeRequest struct {
	Name string                `json:"name" binding:"required"`
	Type exchange.ExchangeType `json:"type" binding:"required"`
}

func (c *Controller) CreateExchange(ctx *gin.Context) {

	var (
		request = CreateExchangeRequest{}
	)

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		logrus.Error("error in binding request ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = c.broker.CreateExchange(ctx, request.Name, request.Type)
	if err != nil {
		logrus.Error("error in creating exchange ", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}
