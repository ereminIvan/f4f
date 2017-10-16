package app

import (
	"log"
	"sync"
	"strings"
	"time"

	"github.com/ereminIvan/fffb/model"
	"github.com/ereminIvan/fffb/service"
)

type application struct {
	config model.Config

	fbService     IFBService
	tgService     ITelegramService
	filterService IFilterService

	sync.Once
}

var (
	fbDumpFilePath string
	tgDumpFilePath string
	configFilePath string
)

func Init(configPath, fbDumpPath, tgDumpPath string) (*application, error) {
	fbDumpFilePath = fbDumpPath
	tgDumpFilePath = tgDumpPath
	configFilePath = configPath

	log.Print("Init Application ...")
	var err error
	a := &application{}

	if a.config, err = a.readConfig(); err != nil {
		return nil, err
	}

	fbDump, err := a.readFBDump()
	if err != nil {
		log.Printf("Cant read FB dump: %v", err)
	}
	tgDump, err := a.readTGDump()
	if err != nil {
		log.Printf("Cant read TG dump: %v", err)
	}

	a.initServices(fbDump, tgDump)

	return a, nil
}

//initServices - init services
func (a *application) initServices(fbDump map[string]struct{}, tgDump map[int64]struct{}) {
	log.Print("Init services ...")
	a.Once.Do(func() {
		a.filterService = service.NewFilterService(a.config.Filter)
		a.fbService = service.NewFBService(a.config.FB, fbDump)
		a.tgService = service.NewTelegramService(a.config.Telegram, tgDump)
	})
}

//Run application
func (a *application) Run(shutdown chan struct{}) {
	log.Print("Run application ...")

	for {
		log.Printf(strings.Repeat("=", 100))
		//used non blocking chanel read for stopping request loop
		select {
		case <-shutdown:
			return
		default:
		}
		newMessages := a.fbService.LatestMessages()

		go func() {
			for _, message := range newMessages {
				a.filterService.SetType(&message)
				if message.Type == model.MessageTypeSpam || message.Type == model.MessageTypeTenant {
					continue
				}
				a.tgService.SendMessage(message)
				time.Sleep(1 * time.Second)
			}
		}()
		log.Printf("Pause before next FB request %d seconds", time.Duration(a.config.FB.Delay))
		time.Sleep(time.Duration(a.config.FB.Delay) * time.Second)
	}

	a.Stop()
}

//Stop graceful stop application
func (a *application) Stop() {
	log.Print("Shutdown awaiting finishing")

	a.dumpData()

	log.Print("Shutdown application ... awaiting finish")
}

//dumpData
func (a *application) dumpData() {
	log.Print("Dumping data ...")
	if err := a.writeFBDump(a.fbService.ReadMessages()); err != nil {
		log.Printf("Dump: Error during dump FB messages %v", err)
	}
	if err := a.writeTGDump(a.tgService.Chats()); err != nil {
		log.Printf("Dump: Error during dump FB messages %v", err)
	}
}
