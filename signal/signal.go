package signal

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	SIGNAL_CLI_DBUS_SERVICE = "org.asamk.Signal"
)

type Signal struct {
	Username string
	Messages chan *Message
}

func NewSignal(username string) *Signal {
	messages := make(chan *Message, 10)
	return &Signal{username, messages}
}

// Listen() establishes a connection to the DBus service and listens for
// incoming Signal messages.
func (s *Signal) Listen() {
	log.Print("Connecting to signal-cli.")
	if err := launchSignalCLI(s.Username); err != nil {
		log.Fatalf("Unable to start signal-cli: %v", err)
	}
	signals := make(chan *dbus.Signal, 10)
	conn, err := connectDBus(signals)
	if err != nil {
		log.Fatalf("Failed to connect to DBus: %v", err)
	}
	defer conn.Close()

	for signal := range signals {
		log.Print(signal)
		// Read receipts are of no interest to this application.
		if signal.Name == "org.asamk.Signal.ReceiptReceived" {
			continue
		}
		message, err := newMessageFromSignal(signal)
		if err != nil {
			log.Printf("Failed to parse new message from signal: %v", err)
			continue
		}
		s.Messages <- message
	}
}

// SendMessage() uses the signal-cli command in dbus mode. This is due to
// the method org.asamk.Signal.sendMessage not working as expected when called
// using dbus.Object.Call().
// TODO: further investigate using org.asamk.Signal.sendMessage method.
func (s *Signal) SendMessage(msg *Message) (err error) {
	args := []string{"--dbus", "send", "-m", msg.Text, msg.Username}
	if len(msg.Attachments) > 0 {
		args = append(args, "-a")
		for _, attachment := range msg.Attachments {
			args = append(args, attachment)
		}
	}
	cmd := exec.Command("signal-cli", args...)
	_, err = cmd.Output()
	if err != nil {
		return
	}
	log.Println("Message sent")
	return
}

// Launches signal-cli on the Session DBus in daemon mode.
func launchSignalCLI(username string) (err error) {
	cmd := exec.Command("signal-cli", "-u", username, "daemon")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return
	}
	log.Print("Started signal-cli on the session bus")
	return
}

// Establish a connection to the org.asamk.Signal DBus interface on the Session Bus.
func connectDBus(signals chan<- *dbus.Signal) (conn *dbus.Conn, err error) {
	conn, err = dbus.SessionBus()
	if err != nil {
		return
	}
	for true {
		if err = verifyConnection(conn); err != nil {
			log.Print(err)
			time.Sleep(1 * time.Second)
		} else {
			log.Print("signal-cli connection success!")
			break
		}
	}
	options := dbus.WithMatchSender(SIGNAL_CLI_DBUS_SERVICE)
	err = conn.AddMatchSignal(options)
	if err != nil {
		return
	}
	conn.Signal(signals)
	return
}

// Verify that org.asamk.Signal is available via the DBus interface.
func verifyConnection(conn *dbus.Conn) (err error) {
	s := []string{}
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		return
	}
	found := false
	for _, v := range s {
		if v == SIGNAL_CLI_DBUS_SERVICE {
			found = true
			break
		}
	}
	if !found {
		err = errors.New("signal-cli connection not found on Dbus.")
	}
	return
}
