package signal

import (
	"errors"
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
)

type Message struct {
	Time        time.Time
	Recipient   string
	Text        string
	Attachments []string
}

func NewMessage(t time.Time, recipient, text string, attachments []string) *Message {
	return &Message{
		Time:        t,
		Recipient:   recipient,
		Text:        text,
		Attachments: attachments,
	}
}

// Helper method that transforms a *dbus.Signal to a *Message.
func newMessageFromSignal(signal *dbus.Signal) (msg *Message, err error) {
	utc, ok := signal.Body[0].(int64)
	if !ok {
		err = errors.New(fmt.Sprintf("failed to convert time to int64, %v\n", signal.Body[0]))
		return
	}
	t := time.Unix(utc, 0)
	recipient, ok := signal.Body[1].(string)
	if !ok {
		err = errors.New(fmt.Sprintf("failed to convert recipient number to string, %v\n", signal.Body[1]))
		return
	}
	text, ok := signal.Body[3].(string)
	if !ok {
		err = errors.New(fmt.Sprintf("failed to convert message text to string, %v\n", signal.Body[3]))
		return
	}
	attachments, ok := signal.Body[4].([]string)
	if !ok {
		err = errors.New(fmt.Sprintf("failed to convert attachment path to slice of strings, %v", signal.Body[4]))
		return
	}
	msg = NewMessage(t, recipient, text, attachments)
	return
}
