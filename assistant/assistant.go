package assistant

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/stayradiated/zwolf-signal/signal"
)

type Assistant struct {
	cli      *signal.Signal
	owner    string
	commands []command
}

func NewAssistant(botNumber, ownerNumber string) *Assistant {
	return &Assistant{
		cli:      signal.NewSignal(botNumber),
		owner:    ownerNumber,
		commands: getAllCommands(),
	}
}

func (a *Assistant) Run() {
	go a.cli.Listen()

	for msg := range a.cli.Messages {
		err := a.validateMessage(msg)
		if err != nil {
			a.errorHandler("Failed to validate message,", err)
			continue
		}
		err = a.executeCommand(msg)
		if err != nil {
			a.errorHandler("Failed to execute command,", err)
		}
	}
}

func (a *Assistant) validateMessage(msg *signal.Message) error {
	if msg.Recipient != a.owner {
		return errors.New(fmt.Sprintf("Message arrived from unknown number %v", msg.Recipient))
	}
	if len(msg.Text) == 0 || string(msg.Text[0]) != "!" {
		return errors.New(fmt.Sprintf("Invalid command format. Must start with !. Message Text: %v", msg.Text))
	}
	return nil
}

func (a *Assistant) executeCommand(msg *signal.Message) (err error) {
	splitMessage := strings.Split(msg.Text, " ")
	args := []string{}
	if len(splitMessage) > 1 {
		args = splitMessage[1:]
	}
	command := command{
		cmd:  splitMessage[0],
		args: args,
	}
	switch command.cmd {
	case MAN:
		err = a.returnManual()
		if err != nil {
			return
		}
	default:
		err = errors.New("Invalid command, type !man to see a list of available commands.")
	}
	return
}

// Handles !man command.
func (a *Assistant) returnManual() (err error) {
	manual := getCommandManual()
	err = a.sendMessage(manual, nil)
	return
}

// Wraps signal.SendMessage.
func (a *Assistant) sendMessage(text string, attachments []string) (err error) {
	msg := signal.NewMessage(time.Now(), a.owner, text, attachments)
	err = a.cli.SendMessage(msg)
	return
}

// A helper method to log the error and send a notification to the Assistant owner.
func (a *Assistant) errorHandler(message string, err error) {
	txt := fmt.Sprintf("%v %v", message, err)
	log.Print(txt)
	// Notify the owner of any errors.
	err = a.sendMessage(txt, nil)
	if err != nil {
		log.Printf("Failed to send message, %v", err)
	}
}

// Copy the src file to dst. Any existing file will be overwritten.
func copy(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}
	err = out.Close()
	return
}
