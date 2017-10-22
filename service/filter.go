package service

import (
	"log"
	"fmt"
	"strings"

	"github.com/iveronanomi/fffb/model"
)

const (
	keywordsSeparator = " "
)

type filterService struct {
	spamKeywords     []string // keywords detection of spam message
	landlordKeywords []string // keywords detection of landlord message
	tenantKeywords   []string // keywords detection of tenant message

	currentText      string
	currentTextSlice []string
}

func NewFilterService(cfg model.Filter) *filterService {
	spam := strings.Split(cfg.KeywordsSpam, keywordsSeparator)
	log.Printf("Filter Service : Spam keywords %+v", spam)
	landlord := strings.Split(cfg.KeywordsLandlord, keywordsSeparator)
	log.Printf("Filter Service : Landlord keywords %+v", landlord)
	tenant := strings.Split(cfg.KeywordsTenant, keywordsSeparator)
	log.Printf("Filter Service : Tenant keywords %+v", tenant)

	return &filterService{
		spamKeywords:     spam,
		landlordKeywords: landlord,
		tenantKeywords:   tenant,
	}
}

func (fs *filterService) SetType(message *model.Message) {
	fs.currentText = strings.ToLower(message.Message)
	fs.currentTextSlice = strings.Fields(fs.currentText)

	landlordIdx := fs.getIndex(fs.landlordKeywords)
	tenantIdx := fs.getIndex(fs.tenantKeywords)
	spamIdx := fs.getIndex(fs.spamKeywords)

	message.AppendDebug(fmt.Sprintf("landlord index: %f\ntenant index: %f\nspam index: %f", landlordIdx, tenantIdx, spamIdx))

	if spamIdx > landlordIdx {
		message.Type = model.MessageTypeSpam
		return
	}
	if tenantIdx > landlordIdx {
		message.Type = model.MessageTypeTenant
		return
	}
	if landlordIdx > 0 {
		message.Type = model.MessageTypeLandlord
		return
	}
	message.Type = model.MessageTypeUnknown
	return
}

func (fs *filterService) getIndex(dic []string) float32 {
	var occurred int

	for _, kw := range dic {
		occurred += strings.Count(fs.currentText, kw)
	}

	if occurred > 0 {
		idx := float32(occurred) / float32(len(fs.currentTextSlice))
		return idx
	}
	return float32(occurred)
}
