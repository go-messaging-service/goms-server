package services

import (
	"../logger"
	"../material"
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strconv"
)

type Message material.Message // just simplify the access to the Message struct

type connectionHandler struct {
	connection      *net.Conn
	RegisterEvent   []func(connectionHandler, []string) // will be fired when a client registeres himself at some topics
	UnregisterEvent []func(connectionHandler, []string) // will be fired when a client un-registeres himself at some topics
}

func (ch *connectionHandler) HandleConnection() {
	if ch.connection == nil {
		logger.Error("Connection not set!")
		return
	}

	ch.waitFor([]string{material.MtRegister}, []func(Message){ch.handleRegistration})
	//TODO implement them:
	//	ch.handleSending()
	//	ch.handleClose()
	//TODO handle logout
}

func (ch *connectionHandler) waitFor(messageTypes []string, handler []func(message Message)) {
	rawMessage, err := bufio.NewReader(*ch.connection).ReadString('\n')

	for err == nil {
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
			if message.MessageType == messageType {
				handler[i](message)
			}
		}

		// read again...
		rawMessage, err = bufio.NewReader(*ch.connection).ReadString('\n')
	}
}

func (ch *connectionHandler) getMessageFromJSON(jsonData string) Message {
	message := Message{}
	json.Unmarshal([]byte(jsonData), &message)
	return message
}

func (ch *connectionHandler) handleRegistration(message Message) {
	logger.Debug(fmt.Sprintf("%#v", message))
	logger.Debug(strconv.Itoa(len(ch.RegisterEvent)))
	for _, event := range ch.RegisterEvent {
		event(*ch, message.Topics)
	}
}

func (ch *connectionHandler) handleSending(message Message) {

}

func (ch *connectionHandler) handleClose(message Message) {

}
