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

type application struct {
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

func Init() (*application, error) {
	a := &application{}
	if err := a.readConfig(); err != nil {
		return nil, err
	}
	a.initServices()
	return a, nil
}

func (a *application) readConfig() error {
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		return errors.New("Error during read config.json: " + err.Error())
	}

	if err := json.Unmarshal(file, &a.config); err != nil {
		return errors.New("Error during unmarshal config.json: " + err.Error() + "\n" + string(file))
	}

	log.Printf("Starting Application with config: %#v", a.config)

	return nil
}

func (a *application) initServices() {
	a.Once.Do(func() {
		a.fbService = service.NewFBService(a.config.FB)
		a.telegramService = service.NewTelegramService(a.config.Telegram)
	})
}

func (a *application) Run() {
	log.Print("Start application ...")
	for {
		newMessages := a.fbService.GetLastFeedMessages()

		for _, message := range newMessages {
			a.telegramService.SendMessage(message)
			time.Sleep(1 * time.Second)
		}

		time.Sleep(time.Duration(a.config.FB.FeedRequestFrequency) * time.Minute)
	}
}

func (a *application) Finish() {
	log.Print("Finalizing application ...")
	a.dumpData()
}

func (a *application) dumpData() {
	log.Print("Dumping data ...")
}
