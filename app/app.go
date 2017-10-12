package app

import (
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

var (
	fbDumpFilePath string
	configFilePath string
)

func Init(configPath, fbDumpPath string) (*application, error) {
	fbDumpFilePath = fbDumpPath
	configFilePath = configPath

	log.Print("Init Application ...")
	var err error
	a := &application{}

	if a.config, err = a.readConfig(); err != nil {
		return nil, err
	}

	fbDump, err := a.readDumps()
	if err != nil {
		log.Printf("Cant read dump: %v", err)
	}

	a.initServices(fbDump)

	return a, nil
}

//initServices - init services
func (a *application) initServices(fbDump map[string]struct{}) {
	log.Print("Init services ...")
	a.Once.Do(func() {
		a.fbService = service.NewFBService(a.config.FB, fbDump)
		a.telegramService = service.NewTelegramService(a.config.Telegram)
	})
}

//Run application
func (a *application) Run() {
	log.Print("Run application ...")
	pause := time.Duration(a.config.FB.FeedRequestFrequency) * time.Minute

	for {
		newMessages := a.fbService.LatestMessages()

		go func() {
			for _, message := range newMessages {
				a.telegramService.SendMessage(message)
				time.Sleep(1 * time.Second)
			}
		}()

		log.Printf("Pause before next FB request %d seconds", pause/time.Second)
		time.Sleep(pause)
	}

	a.Stop()
}

func (a *application) Stop() {
	log.Print("Shutdown application ...")
	a.dumpData()
}

func (a *application) dumpData() {
	log.Print("Dumping data ...")
	if err := a.writeFBDump(a.fbService.ReadMessages()); err != nil {
		log.Printf("Dump: Error during dump FB messages %v", err)
	}
}
