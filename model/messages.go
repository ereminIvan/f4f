package model

import (
	"fmt"
	"log"
	"strings"
)

//Message
type Message struct {
	Id         string
	UpdateTime string
	Message    string
}

//String build message text
func (m *Message) String() string {
	parts := strings.Split(m.Id, "_")
	if len(parts) != 2 {
		log.Printf("Invalid message ID value: %s ", m.Id)
		return m.Message
	}
	link := fmt.Sprintf(`<a href="https://m.facebook.com/groups/%s?view=permalink&id=%s">post link</a>`, parts[0], parts[1])
	return fmt.Sprintf("%s\n%s\n%s", link, m.UpdateTime, m.Message)
}