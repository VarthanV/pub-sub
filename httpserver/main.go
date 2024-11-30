package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitServer(port string, controller *HTTPController) {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "PONG"})
	})

	logrus.Infof("HTTP server listening on %s", port)
	r.Run(fmt.Sprintf(":%s", port))

}
