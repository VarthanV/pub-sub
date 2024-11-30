package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type messageType string

const (
	messageTypeSubscribe messageType = "SUBSCRIBE"
)

type message struct {
	Type    messageType            `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type subscribeRequest struct {
	QName string `mapstructure:"queue_name" json:"queue_name"`
}

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
		var (
			msg = message{}
		)
		err := conn.ReadJSON(&msg)
		if err != nil {
			logrus.Error("error in reading connection ", err)
			break
		}

		switch msg.Type {
		case messageTypeSubscribe:
			var (
				request = subscribeRequest{}
			)

			logrus.Info("Subscription request")
			err := mapstructure.Decode(msg.Payload, &request)
			if err != nil {
				logrus.Error("unable to decode map ", err)
			}

			if request.QName == "" {
			}
			err = c.broker.Subscribe(ctx, request.QName, conn)
			if err != nil {
				logrus.Error("error in subscribing to queue ", err)
			}
		}
	}
}
