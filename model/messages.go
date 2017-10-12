package model

import (
	"fmt"
	"log"
	"strings"
)

type MessageType int

const (
	MessageTypeUnknown  MessageType = iota
	MessageTypeSpam
	MessageTypeLandlord
	MessageTypeTenant
)

var messageTypes = map[MessageType]string{
	MessageTypeUnknown:  "unknown",
	MessageTypeSpam:     "spam",
	MessageTypeLandlord: "landlord",
	MessageTypeTenant:   "tenant",
}

func (t MessageType) String() string {
	return messageTypes[t]
}

//Message
type Message struct {
	Id         string
	UpdateTime string
	Message    string
	Type       MessageType
}

//String build message text
func (m *Message) String() string {
	parts := strings.Split(m.Id, "_")
	if len(parts) != 2 {
		log.Printf("Invalid message ID value: %s ", m.Id)
		return m.Message
	}
	link := fmt.Sprintf(`<a href="https://m.facebook.com/groups/%s?view=permalink&id=%s">post link</a>`, parts[0], parts[1])
	return fmt.Sprintf("%s %s\n%s\n%s", link, m.Type.String(), m.UpdateTime, m.Message)
}