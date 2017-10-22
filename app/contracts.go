package app

import "github.com/iveronanomi/fffb/model"

type IFBService interface {
	LatestMessages() []model.Message
	ReadMessages() []string
}

type ITelegramService interface {
	SendMessage(message model.Message)
	Chats() map[int64]struct{}
}

type IFilterService interface {
	SetType(message *model.Message)
}