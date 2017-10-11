package service

import (
	"log"

	api "github.com/huandu/facebook"

	"github.com/ereminIvan/fffb/model"
)

type fbService struct {
	config  model.FBConfig
	session *api.Session

	radMessages map[string]map[string]string
}

func NewFBService(cfg model.FBConfig) *fbService {
	radMessages := make(map[string]map[string]string)

	return &fbService{
		config:      cfg,
		radMessages: radMessages,
	}
}

func (s *fbService) GetLastFeedMessages() []model.Message {
	// create a global App var to hold app id and secret.
	var globalApp = api.New(s.config.AppId, s.config.AppSecret)

	// if there is another way to get decoded access token,
	// creates a session directly with the token.
	if s.session == nil {
		s.session = globalApp.Session(globalApp.AppAccessToken())
	}

	if s.config.DebugEnabled {
		s.session.SetDebug(api.DebugMode(s.config.DebugMode))
	}

	// validate access token. err is nil if token is valid.
	if err := s.session.Validate(); err != nil {
		log.Printf("FB: Error ocured during fb session validation: %v", err)
	}

	// use session to send api request with access token.
	res, _ := s.session.Get(s.config.FeedURL, nil)

	r := []map[string]string{}
	res.DecodeField("data", &r)

	return s.processMessages(r)
}

func (s *fbService) processMessages(m []map[string]string) []model.Message {
	result := []model.Message{}

	log.Printf("FB: Count of old messages is : %d", len(s.radMessages))

	for _, item := range m {
		if _, ok := s.radMessages[item["id"]]; ok {
			continue
		}
		s.radMessages[item["id"]] = item
		msg := model.Message{
			Message:    item["message"],
			UpdateTime: item["updated_time"],
			Id:         item["id"],
		}
		result = append(result, msg)
	}
	log.Printf("FB: Count of new messages is: %d", len(result))

	return result
}
