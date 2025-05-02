package domain

import (
	"VK_posts/internal/logger"
	"VK_posts/internal/models"
	"github.com/labstack/echo"
	"github.com/olahol/melody"
	"net/http"
	"strings"
	"sync"
)

var (
	mains       = make(map[string]*melody.Session)
	messageHash = make(chan models.Message, 20)
)

func MessageConvert(msg string) models.Message {
	// message like a "From/to/content/time
	message := strings.Split(msg, "/")
	return models.Message{From: message[0], To: message[1], Content: message[2], Timestamp: message[3]}
}

type MessengerDomain struct {
	m                        *melody.Melody
	mu                       sync.Mutex
	PostgresMessengerHandler PostgresMessengerHandler
}
type PostgresMessengerHandler interface {
	GetUserServer(userid string) (string, error)
	MessageSave(message models.Message) error
}

func NewMessengerDomain(handler PostgresMessengerHandler) *MessengerDomain {
	m := melody.New()
	var mu sync.Mutex
	return &MessengerDomain{m: m, mu: mu, PostgresMessengerHandler: handler}
}
func (d *MessengerDomain) MessengerLogic(c echo.Context) error {
	op := "MessengerLogic"
	log := logger.GetLogger()
	log.With("op", op)
	err := d.m.HandleRequest(c.Response(), c.Request())
	if err != nil {
		log.Error("Problem to turn websocket", "err", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	go func() {
		for message := range messageHash {
			err = d.PostgresMessengerHandler.MessageSave(message)
			if err != nil {
				log.Error("Problem to message save ", "err", err.Error())
				messageHash <- message
			}
		}
	}()
	d.m.HandleConnect(func(s *melody.Session) {
		userID := c.Get("userID").(string)
		main := c.Get("Main").(string)
		if main == "true" {
			d.mu.Lock()
			defer d.mu.Unlock()
			s.Set("userID", userID)
			mains[userID] = s
		}
	})
	d.m.HandleMessage(func(s *melody.Session, msg []byte) {
		message := MessageConvert(string(msg))
		err = d.PostgresMessengerHandler.MessageSave(message)
		if err != nil {
			messageHash <- message
		}
		client, ok := mains[message.To]
		if !ok {
			err = client.Write(msg)
			if err != nil {
				s.Close()
			}
		}
		status := "ok"
		err = s.Write([]byte(status))
		if err != nil {
			s.Close()
		}
	})
	d.m.HandleDisconnect(func(s *melody.Session) {
		d.mu.Lock()
		defer d.mu.Unlock()
		value, status := s.Get("userID")
		if status {
			err = mains[value.(string)].Close()
			if err != nil {
				log.Error("Problem to close main session", "err", err.Error())
				//Todo если ошибка что не существует то пропускаем
			}
		}
	})
	return nil
}
