package service

import (
	"log"

	api "gopkg.in/telegram-bot-api.v4"

	"github.com/iveronanomi/fffb/model"
)

type tgService struct {
	config model.TelegramConfig
	chats  map[int64]struct{} // list of chat ids
	bot    *api.BotAPI
}

type IFBService interface {
	SendMessage(message model.Message)
}

func NewTelegramService(cfg model.TelegramConfig, chats map[int64]struct{}) *tgService {
	var err error
	s := &tgService{
		config: cfg,
		chats:  chats,
	}
	s.bot, err = api.NewBotAPI(cfg.Token)

	if err != nil {
		log.Panic(err)
	}

	s.bot.Debug = cfg.DebugEnabled

	log.Printf("Telegram: Authorized on account %s", s.bot.Self.UserName)

	return s
}

func (s *tgService) SendMessage(message model.Message) {
	log.Print("Telegram: Sending message")
	u := api.UpdateConfig{Limit: 100, Offset: 1, Timeout: 60}

	updates, err := s.bot.GetUpdates(u)

	if err != nil {
		panic(err)
	}

	log.Print("Telegram: SendMessage [read updates]")
	//read all updates
	for _, update := range updates {
		//append new user chats if not exist | don't check existence because chat id could be new
		s.chats[update.Message.Chat.ID] = struct{}{}
		//log.Printf("Telegram: SendMessage [%d][%s] %s", update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text)
	}
	log.Printf("Telegram: SendMessage Chats:: %+v", s.chats)
	//send to all subscribers new message
	for chatID := range s.chats {
		text := message.String()
		log.Printf("Telegram: SendMessage [chat]:%d | [text len]:%d", chatID, len(text))
		s.bot.Send(NewMessage(chatID, text))
	}
}

func NewMessage(chatID int64, message string) api.MessageConfig {
	return api.MessageConfig{
		BaseChat: api.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		Text:                  message,
		ParseMode:             "HTML",
		DisableWebPagePreview: false,
	}
}

func (s *tgService) Chats() map[int64]struct{} {
	return s.chats
}

//todo remove it
func (s *tgService) UpdateChanel() {
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

		msg := NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		s.bot.Send(msg)
	}
}
