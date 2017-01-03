package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"goMS/src/material"
	"goMS/src/technical/common"
	"goMS/src/technical/services/logger"
	"math"
	"net"
)

type Message material.Message // just simplify the access to the Message struct

type connectionHandler struct {
	connection       *net.Conn
	connectionClosed bool
	registeredTopics []string
	RegisterEvent    []func(connectionHandler, []string) // will be fired when a client registeres himself at some topics
	UnregisterEvent  []func(connectionHandler, []string) // will be fired when a client un-registeres himself at some topics
	SendEvent        []func(connectionHandler, []string, string)
}

func (ch *connectionHandler) Init(connection *net.Conn) {
	ch.connection = connection
}

func (ch *connectionHandler) HandleConnection() {
	// Not initialized
	if ch.connection == nil {
		logger.Error("Connection not set!")
		return
	}

	// at first only a registration message is allowed
	ch.waitFor(
		[]string{material.MT_REGISTER},
		[]func(Message){ch.handleRegistration})

	// Now a arbitrary amount of registration, logout, close and send messages is allowed
	for !ch.connectionClosed {
		ch.waitFor(
			[]string{material.MT_REGISTER,
				material.MT_LOGOUT,
				material.MT_CLOSE,
				material.MT_SEND},
			[]func(Message){ch.handleRegistration,
				ch.handleLogout,
				ch.handleClose,
				ch.handleSending})

		// When the connection is closed, exit the loop and we're done
		if ch.connectionClosed {
			break
		}
	}
}

func (ch *connectionHandler) waitFor(messageTypes []string, handler []func(message Message)) {
	rawMessage, err := bufio.NewReader(*ch.connection).ReadString('\n')

	if err == nil {
		// the length of the message that should be printed
		maxOutputLength := int(math.Min(float64(len(rawMessage))-1, 30))
		output := rawMessage[:maxOutputLength]
		if 30 < len(rawMessage)-1 {
			output += " [...]"
		}
		logger.Info(output)

		// JSON to Message-struct
		message := ch.getMessageFromJSON(rawMessage)

		// check type
		for i := 0; i < len(messageTypes); i++ {
			messageType := messageTypes[i]
			logger.Info("Check " + messageType + " type")

			if message.MessageType == messageType {
				logger.Info("Handle " + messageType + " type")
				handler[i](message)
				break
			}
		}
	} else {
		logger.Info("The connection will be closed. Reason: " + err.Error())
		ch.exit()
	}
}

func (ch *connectionHandler) getMessageFromJSON(jsonData string) Message {
	message := Message{}
	json.Unmarshal([]byte(jsonData), &message)
	return message
}

func (ch *connectionHandler) handleRegistration(message Message) {
	logger.Debug("Register to topics " + fmt.Sprintf("%#v", message.Topics))

	for _, event := range ch.RegisterEvent {
		event(*ch, message.Topics)
	}

	for _, topic := range message.Topics {
		if !common.ContainsString(ch.registeredTopics, topic) {
			ch.registeredTopics = append(ch.registeredTopics, topic)
		}
	}
}

func (ch *connectionHandler) handleSending(message Message) {
	logger.Debug(fmt.Sprintf("Send message to topics %#v", message.Topics))

	for _, event := range ch.SendEvent {
		event(*ch, message.Topics, message.Data)
	}
}

func (ch *connectionHandler) handleLogout(message Message) {
	logger.Debug(fmt.Sprintf("Unsubscribe from topics %#v", message.Topics))
	ch.logout(message.Topics)
}

func (ch *connectionHandler) handleClose(message Message) {
	ch.exit()
}

func (ch *connectionHandler) exit() {
	logger.Debug("Unsubscribe from all topics")
	ch.logout(ch.registeredTopics)

	logger.Debug("Close connection")
	(*ch.connection).Close()
	ch.connectionClosed = true
}

func (ch *connectionHandler) logout(topics []string) {
	for _, event := range ch.UnregisterEvent {
		event(*ch, topics)
	}
}
