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

const MAX_PRINTING_LENGTH int = 30

// Init initializes the handler with the given connection.
func (ch *connectionHandler) Init(connection *net.Conn) {
	ch.connection = connection
}

// HandleConnection starts a routine to handle registration and sending messages.
// This will run until the client logs out, so run this in a goroutine.
func (ch *connectionHandler) HandleConnection() {
	// Not initialized
	if ch.connection == nil {
		logger.Fatal("Connection not set!")
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

// waitFor wats until on of the given message types arrived.
// The i-th argument in the messageTypes array must match to the i-th argument in the handler array.
func (ch *connectionHandler) waitFor(messageTypes []string, handler []func(message Message)) {

	// Check if the arrays match and error/fatal here
	if len(messageTypes) != len(handler) {
		if len(messageTypes) > len(handler) {
			// Fatal here to prevent a "slice bounds out of range" error during runtime
			logger.Fatal("There're more defined message types then functions mapped to them.")
		} else {
			logger.Error("There're more defined functions then message types here. Some message types might not be covered. Fix that!")
		}
	}

	rawMessage, err := bufio.NewReader(*ch.connection).ReadString('\n')

	if err == nil {
		// the length of the message that should be printed
		maxOutputLength := int(math.Min(float64(len(rawMessage))-1, float64(MAX_PRINTING_LENGTH)))
		output := rawMessage[:maxOutputLength]
		if MAX_PRINTING_LENGTH < len(rawMessage)-1 {
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

// getMessageFromJSON converts the given json-data into a message object.
func (ch *connectionHandler) getMessageFromJSON(jsonData string) Message {
	message := Message{}
	json.Unmarshal([]byte(jsonData), &message)
	return message
}

// handleRegistration registeres this connection to the topics specified in the message.
func (ch *connectionHandler) handleRegistration(message Message) {
	logger.Debug("Register to topics " + fmt.Sprintf("%#v", message.Topics))

	for _, event := range ch.RegisterEvent {
		event(*ch, message.Topics)
	}

	for _, topic := range message.Topics {
		if !technicalCommon.ContainsString(ch.registeredTopics, topic) {
			ch.registeredTopics = append(ch.registeredTopics, topic)
		}
	}
}

// handleSending send the given message to all clients interested in the topics specified in the message.
func (ch *connectionHandler) handleSending(message Message) {
	logger.Debug(fmt.Sprintf("Send message to topics %#v", message.Topics))

	for _, event := range ch.SendEvent {
		event(*ch, message.Topics, message.Data)
	}
}

// handleLogout logs the client out.
func (ch *connectionHandler) handleLogout(message Message) {
	logger.Debug(fmt.Sprintf("Unsubscribe from topics %#v", message.Topics))
	ch.logout(message.Topics)
}

// handleClose logs the client out from all topics and closes the connection.
func (ch *connectionHandler) handleClose(message Message) {
	ch.exit()
}

// exit logs the client out from all topics and closes the connection.
func (ch *connectionHandler) exit() {
	logger.Debug("Unsubscribe from all topics")
	ch.logout(ch.registeredTopics)

	logger.Debug("Close connection")
	(*ch.connection).Close()
	ch.connectionClosed = true
}

// logout will logs the client out from the given topics.
func (ch *connectionHandler) logout(topics []string) {
	for _, event := range ch.UnregisterEvent {
		event(*ch, topics)
	}
	//TODO also remove topics from ch.registeredTopics lol
}
