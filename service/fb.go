package service

import (
	"log"

	api "github.com/huandu/facebook"

	"github.com/ereminIvan/fffb/model"
)

var globalApp *api.App

type fbService struct {
	config  model.FBConfig
	session *api.Session

	readMessages map[string]struct{} // list of read messages
}

func NewFBService(cfg model.FBConfig, dumpMessages map[string]struct{}) *fbService {
	if dumpMessages == nil {
		dumpMessages = make(map[string]struct{})
	}
	return &fbService{
		config:       cfg,
		readMessages: dumpMessages,
	}
}

//GetLastFeedMessages - get latest unread feed messages
func (s *fbService) LatestMessages() []model.Message {
	//log.Print("FB: Getting new messages")
	// create a global App var to hold app id and secret.
	if globalApp == nil {
		globalApp = api.New(s.config.AppId, s.config.AppSecret)
	}
	// if there is another way to get decoded access token,
	// creates a session directly with the token.
	if s.session == nil {
		s.session = globalApp.Session(globalApp.AppAccessToken())
	}

	if s.config.DebugEnabled {
		s.session.SetDebug(api.DebugMode(s.config.DebugMode))
	}

	// validate access token. err is nil if token is valid.
	//if err := s.session.Validate(); err != nil {
	//	log.Printf("FB: Session validation error: %v", err)
	//}

	// use session to send api request with access token.
	res, _ := s.session.Get(s.config.FeedURL, nil)

	r := []map[string]string{}
	res.DecodeField("data", &r)

	messages := s.processMessages(r)

	return messages
}

//processMessages - get latest unread messages list form message pull
func (s *fbService) processMessages(m []map[string]string) []model.Message {
	result := []model.Message{}

	//log.Printf("FB: Count of old messages is: %d", len(s.readMessages))

	for _, item := range m {
		if _, ok := s.readMessages[item["id"]]; ok {
			continue
		}

		//add message id to reade messages
		s.readMessages[item["id"]] = struct{}{}
		result = append(result, model.Message{
			Message:    item["message"],
			UpdateTime: item["updated_time"],
			Id:         item["id"],
		})
	}

	//reverse for sending
	reverse := make([]model.Message, len(result))
	for idx, m := range result {
		reverse[len(result)-1-idx] = m
	}
	//todo check if params has `limit` parameter
	//if len(result) >= int(s.config.FeedLimit) {
	//	result = result[:int(s.config.FeedLimit) - 1]
	//}

	log.Printf("FB: Count of new messages is: %d", len(reverse))

	return reverse
}

//ReadMessages - receive all read messages ids
func (s *fbService) ReadMessages() []string {
	result := make([]string, 0, len(s.readMessages))
	for id := range s.readMessages {
		result = append(result, id)
	}
	return result
}
