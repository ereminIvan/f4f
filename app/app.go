package app

import (
	"io/ioutil"
	"encoding/json"
	"errors"

	"github.com/ereminIvan/fffb/model"
	"github.com/ereminIvan/fffb/service/fb"
)

type Application struct {
	Config model.Config
}

func Init() (*Application, error) {
	a := &Application{}
	if err := a.readConfig(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Application) readConfig() error {
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		return errors.New("Error during read config.json: " + err.Error())
	}

	if err := json.Unmarshal(file, &a.Config); err != nil {
		return errors.New("Error during unmarshal config.json: " + err.Error())
	}

	return nil
}

func (a *Application) Run() {
	fbService := fb.NewService(a.Config.FB)
	fbService.GetFeed()
}
