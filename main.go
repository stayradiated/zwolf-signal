package main

import (
	"log"
	"os"

	"github.com/stayradiated/zwolf-signal/assistant"
)

func main() {
	botNumber := os.Getenv("BOT_NUMBER")
	if botNumber == "" {
		log.Fatal("The BOT_NUMBER env var is required.")
	}
	ownerNumber := os.Getenv("OWNER_NUMBER")
	if ownerNumber == "" {
		log.Fatal("The OWNER_NUMBER env var is required.")
	}
	assistant := assistant.NewAssistant(botNumber, ownerNumber)
	assistant.Run()
}
