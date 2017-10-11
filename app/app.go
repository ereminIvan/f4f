package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/ereminIvan/fffb/model"
	"github.com/ereminIvan/fffb/service"
)

type Application struct {
	config model.Config

	fbService       IFBService
	telegramService ITelegramService

	sync.Once
}

type IFBService interface {
	GetLastFeedMessages() []model.Message
}

type ITelegramService interface {
	SendMessage(message model.Message)
}

func Init() (*Application, error) {
	a := &Application{}
	if err := a.readConfig(); err != nil {
		return nil, err
	}
	a.initServices()
	return a, nil
}

func (a *Application) readConfig() error {
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		return errors.New("Error during read config.json: " + err.Error())
	}

	if err := json.Unmarshal(file, &a.config); err != nil {
		return errors.New("Error during unmarshal config.json: " + err.Error())
	}

	return nil
}

func (a *Application) initServices() {
	a.Once.Do(func() {
		a.fbService = service.NewFBService(a.config.FB)
		a.telegramService = service.NewTelegramService(a.config.Telegram)
	})
}

func (a *Application) Run() {
	for {
		newMessages := a.fbService.GetLastFeedMessages()

		for _, message := range newMessages {
			a.telegramService.SendMessage(message)
			log.Printf("DEBUG: After telegram call")
		}

		//time.Sleep(time.Duration(a.config.FB.FeedRequestFrequency) * time.Minute)
		time.Sleep(2 * time.Second)
	}
}
