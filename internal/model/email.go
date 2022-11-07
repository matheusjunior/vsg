package model

import (
	"fmt"
)

type Email struct {
	From string
	To   []string
	Body string
}

func (e *Email) Message() []byte {
	s := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Test mail\r\n\r\n"+
		"%s\r\n", e.From, e.To, e.Body)
	return []byte(s)
}
