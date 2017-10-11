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

	s.bot.Debug = cfg.DebugEnabled

	log.Printf("Telegram: Authorized on account %s", s.bot.Self.UserName)

	return s
}

func (s *telegramService) SendMessage(message model.Message) {
	u := api.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.bot.GetUpdates(u)

	if err != nil {
		panic(err)
	}

	for _, update := range updates {
		log.Printf("Telegram: [%d][%s] %s", update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text)
		if update.Message == nil {
			continue
		}
		msg := api.NewMessage(
			update.Message.Chat.ID,
			//fmt.Sprintf("Id: %s\nTime: %s\nMessage: %s\n", message.Id, message.UpdateTime, message.Message),
			message.Message,
		)
		s.bot.Send(msg)
	}
}

func (s *telegramService) UpdateChanel() {
	//todo example
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

		log.Printf("Telegram: [%s][%s] %s", update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text)

		msg := api.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		s.bot.Send(msg)
	}
}