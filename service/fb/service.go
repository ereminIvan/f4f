package fb

import (
	"fmt"

	api "github.com/huandu/facebook"

	"github.com/ereminIvan/fffb/model"
)

type service struct {
	config  model.FBConfig
	session *api.Session

	cachedFeed map[string]map[string]string
}

func NewService(cfg model.FBConfig) *service {
	cached := make(map[string]map[string]string)

	return &service{
		config:     cfg,
		cachedFeed: cached,
	}
}

func (s *service) GetFeed() {
	// create a global App var to hold app id and secret.
	var globalApp = api.New(s.config.AppId, s.config.AppSecret)

	// if there is another way to get decoded access token,
	// creates a session directly with the token.
	if s.session == nil {
		s.session = globalApp.Session(globalApp.AppAccessToken())
	}

	// validate access token. err is nil if token is valid.
	if err := s.session.Validate(); err != nil {
		fmt.Println(err)
	}

	// use session to send api request with access token.
	res, _ := s.session.Get(s.config.FeedURL, nil)

	r := []map[string]string{}
	res.DecodeField("data", &r)

	newMessages := s.getNewMessages(r)

	for _, v := range newMessages {
		fmt.Println("--------------------------------------------")
		fmt.Println(v["id"])
		fmt.Println(v["updated_time"])
		fmt.Println(v["message"])
	}
}

func (s *service) getNewMessages(m []map[string]string) []map[string]string {
	result := []map[string]string{}
	for _, item := range m {
		if _, ok := s.cachedFeed[item["id"]]; ok {
			continue
		}
		s.cachedFeed[item["id"]] = item
		result = append(result, item)
	}
	return result
}
