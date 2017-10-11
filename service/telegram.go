package service

import (
	"log"

	api "gopkg.in/telegram-bot-api.v4"

	"github.com/ereminIvan/fffb/model"
)

type telegramService struct {
	config model.TelegramConfig
	bot    *api.BotAPI
}

type IFBService interface {
	SendMessage(message model.Message)
}

func NewTelegramService(cfg model.TelegramConfig) *telegramService {
	var err error
	s := &telegramService{
		config: cfg,
	}
	s.bot, err = api.NewBotAPI(cfg.Token)

	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Telegram: Authorized on account %s", s.bot.Self.UserName)

	return s
}

func (s *telegramService) SendMessage(message model.Message) {
	u := api.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.bot.GetUpdatesChan(u)

	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("Telegram: [%s] %s", update.Message.From.UserName, update.Message.Text)

		//msg := api.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		msg := api.NewMessage(update.Message.Chat.ID, message.Message)
		s.bot.Send(msg)
		log.Printf("Telegram: loop by updates executed")
	}
}
