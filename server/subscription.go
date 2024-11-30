package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	// allow all origins as of now
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Controller) HandleSubscription(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Error("error in upgrading to websocket connection ", err)
		ctx.Error(errors.New("error in upgrading ws connection"))
		return
	}

	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		logrus.Info("Received: " + string(message))

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}
}
