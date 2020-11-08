package main

import (
	"log"
	"os"

	"github.com/stayradiated/zwolf-signal/service"
)

func main() {
	amqpAddress := os.Getenv("AMQP_ADDRESS")
	if amqpAddress == "" {
		log.Fatal("The AMQP_ADDRESS env var is required.")
	}

	botNumber := os.Getenv("BOT_NUMBER")
	if botNumber == "" {
		log.Fatal("The BOT_NUMBER env var is required.")
	}

	ownerNumber := os.Getenv("OWNER_NUMBER")
	if ownerNumber == "" {
		log.Fatal("The OWNER_NUMBER env var is required.")
	}

	service := service.NewService(botNumber, ownerNumber, amqpAddress)
	service.Run()
}
