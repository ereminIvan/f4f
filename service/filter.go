package service

import (
	"log"
	"strings"

	"github.com/ereminIvan/fffb/model"
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

func (fs *filterService) GetType(text string) model.MessageType {
	fs.currentText = strings.ToLower(text)
	fs.currentTextSlice = strings.Fields(fs.currentText)

	landlordIdx := fs.getIndex(fs.landlordKeywords)
	spamIdx := fs.getIndex(fs.spamKeywords)
	tenantIdx := fs.getIndex(fs.tenantKeywords)

	log.Printf("Filter service: landlord message type with index %f > 0", landlordIdx)
	log.Printf("Filter service: tenant message type with index %f > 0.09", tenantIdx)
	log.Printf("Filter service: spam message type with index %f > 0.09", spamIdx)

	if spamIdx > landlordIdx {
		return model.MessageTypeSpam
	}
	if tenantIdx > landlordIdx {
		return model.MessageTypeTenant
	}
	if landlordIdx > 0 {
		return model.MessageTypeLandlord
	}

	log.Print("Filter service: unknown message type with")
	return model.MessageTypeUnknown
}

func (fs *filterService) getIndex(dic []string) float32 {
	var occurred int

	for _, kw := range dic {
		occurred += strings.Count(fs.currentText, kw)
	}

	log.Printf("Filter service [occured: %d]", occurred)

	if occurred > 0 {
		idx := float32(occurred) / float32(len(fs.currentTextSlice))
		log.Printf("float32(occurred)[%d] / float32(len(fs.currentTextSlice))[%d] = %f", occurred, len(fs.currentTextSlice), idx)
		return idx
	}

	return float32(occurred)
}
