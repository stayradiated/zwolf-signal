package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/stayradiated/zwolf-signal/signal"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

var (
	consumeTopic = "events"
	publishTopic = "events-processed"

	logger = watermill.NewStdLogger(
		true,  // debug
		false, // trace
	)
)

type Service struct {
	cli       *signal.Signal
	owner     string
	publisher message.Publisher
}

func NewService(botNumber, ownerNumber, amqpAddress string) *Service {
	queueConfig := amqp.NewDurableQueueConfig(amqpAddress)

	publisher, err := amqp.NewPublisher(queueConfig, logger)
	if err != nil {
		panic(err)
	}

	return &Service{
		cli:       signal.NewSignal(botNumber),
		owner:     ownerNumber,
		publisher: publisher,
	}
}

func (a *Service) Run() {
	go a.cli.Listen()

	for msg := range a.cli.Messages {
		err := a.validateMessage(msg)
		if err != nil {
			a.errorHandler("Failed to validate message,", err)
			continue
		}

		payload, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		a.publisher.Publish(consumeTopic, message.NewMessage(
			watermill.NewUUID(),
			payload,
		))
	}
}

func (a *Service) validateMessage(msg *signal.Message) error {
	if msg.Username != a.owner {
		return errors.New(fmt.Sprintf("Message arrived from unknown number %v", msg.Username))
	}
	if len(msg.Text) == 0 {
		return errors.New(fmt.Sprintf("Message arrived without any text", msg.Text))
	}
	return nil
}

// Wraps signal.SendMessage.
func (a *Service) sendMessage(text string, attachments []string) (err error) {
	msg := signal.NewMessage(time.Now(), a.owner, text, attachments)
	err = a.cli.SendMessage(msg)
	return
}

// A helper method to log the error and send a notification to the Service owner.
func (a *Service) errorHandler(message string, err error) {
	txt := fmt.Sprintf("%v %v", message, err)
	log.Print(txt)
	// Notify the owner of any errors.
	err = a.sendMessage(txt, nil)
	if err != nil {
		log.Printf("Failed to send message, %v", err)
	}
}
