package main

import (
	"./logger"
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
)

type connectionHandler struct {
	connection *net.Conn
}

func (ch *connectionHandler) HandleConnection() {
	if ch.connection == nil {
		logger.Err("Connection not set!")
		return
	}

	//TODO implement them:
	ch.waitFor([]string{mtRegister}, []func(Message){ch.handleRegistration})
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
		message := Message{}
		json.Unmarshal([]byte(rawMessage), &message)

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

func (ch *connectionHandler) handleRegistration(message Message) {
	logger.Debug(fmt.Sprintf("%#v", message))
}

func (ch *connectionHandler) handleSending(message Message) {

}

func (ch *connectionHandler) handleClose(message Message) {

}
