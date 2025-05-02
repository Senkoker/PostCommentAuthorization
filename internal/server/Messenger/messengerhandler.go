package Messenger

import (
	"VK_posts/internal/models"
	"github.com/labstack/echo"
	"github.com/olahol/melody"
	"net/http"
	"strings"
)

var (
	mains       = make(map[string]*melody.Session)
	messageHash = make(chan models.Message, 20)
)

func MessageConvert(msg string) models.Message {
	message := strings.Split(msg, "/")
	return models.Message{From: message[0], To: message[1], Content: message[2], Timestamp: message[3]}
}

type Messenger interface {
	MessengerLogic(echo.Context) error
}
type ServerLog interface {
	GetRecipientServerInfo(recipientID string) (string, error)
}

func MessengerHandler(messenger Messenger) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := messenger.MessengerLogic(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
func ServerInfo(log ServerLog) echo.HandlerFunc {
	return func(c echo.Context) error {
		recipientID := c.Response().Header().Get("recipientID")
		serverID, err := log.GetRecipientServerInfo(recipientID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, serverID)
	}
}
