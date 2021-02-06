package service

import (
	"context"
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
	sendMessagesTopic     = "send-messages"
	receivedMessagesTopic = "received-messages"

	logger = watermill.NewStdLogger(
		true,  // debug
		false, // trace
	)
)

type Service struct {
	cli        *signal.Signal
	owner      string
	publisher  message.Publisher
	subscriber message.Subscriber
}

func NewService(botNumber, ownerNumber, amqpAddress string) *Service {
	queueConfig := amqp.NewDurableQueueConfig(amqpAddress)

	publisher, err := amqp.NewPublisher(queueConfig, logger)
	if err != nil {
		panic(err)
	}

	subscriber, err := amqp.NewSubscriber(queueConfig, logger)
	if err != nil {
		panic(err)
	}

	return &Service{
		cli:        signal.NewSignal(botNumber),
		owner:      ownerNumber,
		publisher:  publisher,
		subscriber: subscriber,
	}
}

func (s *Service) Run() {
	go s.cli.Listen()
	go s.subscribeToAmqp()

	for msg := range s.cli.Messages {
		err := s.validateMessage(msg)
		if err != nil {
			s.errorHandler("Failed to validate message,", err)
			continue
		}

		fmt.Println(msg)

		payload, err := json.Marshal(msg)
		if err != nil {
			s.errorHandler("Could not marshal message into JSON", err)
		}

		s.publisher.Publish(receivedMessagesTopic, message.NewMessage(
			watermill.NewUUID(),
			payload,
		))
	}
}

func (s *Service) subscribeToAmqp() {
	messages, err := s.subscriber.Subscribe(context.Background(), sendMessagesTopic)
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		fmt.Printf("received message: %s, payload: %s\n", msg.UUID, string(msg.Payload))
		msg.Ack()

		smsg := signal.Message{}
		err := json.Unmarshal(msg.Payload, &smsg)
		if err != nil {
			log.Print(err)
			continue
		}

		s.sendMessage(smsg.Text, smsg.Attachments)
	}
}

func (s *Service) validateMessage(msg *signal.Message) error {
	if msg.Username != s.owner {
		return errors.New(fmt.Sprintf("Message arrived from unknown number %v", msg.Username))
	}
	if len(msg.Text) == 0 {
		return errors.New(fmt.Sprintf("Message arrived without any text", msg.Text))
	}
	return nil
}

// Wraps signal.SendMessage.
func (s *Service) sendMessage(text string, attachments []string) (err error) {
	msg := signal.NewMessage(time.Now(), s.owner, text, attachments)
	err = s.cli.SendMessage(msg)
	return
}

// A helper method to log the error and send a notification to the Service owner.
func (s *Service) errorHandler(message string, err error) {
	txt := fmt.Sprintf("%v %v", message, err)
	log.Print(txt)
	// Notify the owner of any errors.
	err = s.sendMessage(txt, nil)
	if err != nil {
		log.Printf("Failed to send message, %v", err)
	}
}
