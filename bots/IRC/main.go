package main

import (
	"context"
	"github.com/carverauto/threadr/bots/IRC/pkg/adapters/broker"
	irc "github.com/carverauto/threadr/bots/IRC/pkg/adapters/messages"
	"github.com/carverauto/threadr/bots/IRC/pkg/common"
	pm "github.com/carverauto/threadr/bots/IRC/pkg/ports"
	"log"
	"time"
)

// MessageBroker is a type that holds information about the message broker
type MessageBroker struct {
	URL       string
	Subject   string
	Stream    string
	CmdStream string
}

func main() {
	natsURL := "nats://nats.nats.svc.cluster.local:4222"
	subject := "irc"
	stream := "message_processing"
	cmds_subject := "commands"
	cmds_stream := "commands_processing"

	cloudEventsHandler, err := broker.NewCloudEventsNATSHandler(natsURL, subject, stream)
	if err != nil {
		log.Fatalf("Failed to create CloudEvents handler: %s", err)
	}

	commandsHandler, err := broker.NewCloudEventsNATSHandler(natsURL, cmds_subject, cmds_stream)
	if err != nil {
		log.Fatalf("Failed to create CloudEvents handler: %s", err)
	}

	var ircAdapter pm.MessageAdapter = irc.NewIRCAdapter()
	if err := ircAdapter.Connect(commandsHandler); err != nil {
		log.Fatal("Failed to connect to IRC:", err)
	}

	// start a counter for received message_processing
	msgCounter := 0
	ircAdapter.Listen(func(ircMsg common.IRCMessage) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Create a CloudEvent
		ce := broker.Message{
			Sequence:  msgCounter,
			Message:   ircMsg.Message,
			Nick:      ircMsg.Nick,
			Channel:   ircMsg.Channel,
			Platform:  "IRC",
			Timestamp: time.Now(),
		}

		err := cloudEventsHandler.PublishEvent(ctx, ce)
		if err != nil {
			log.Printf("Failed to send CloudEvent: %v", err)
		} else {
			log.Printf("Sent CloudEvent for message [%d]", msgCounter)
		}
		msgCounter++
	})

	// Keep the application running
	select {}
}
