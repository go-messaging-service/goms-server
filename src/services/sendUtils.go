package services

import (
	"encoding/json"
	"goMS/src/material"
	domain "goMS/src/material"
	"goMS/src/technical/services/logger"
	"net"
)

// sendErrorMessage sends the given error data as an error message to the given client.
func sendErrorMessage(conn *net.Conn, errorCode, errorData string) {

	errorMessage := ErrorMessage{
		GenerallMessage: material.GenerallMessage{
			MessageType: material.MT_ERROR,
		},
		ErrorCode: errorCode,
		Error:     errorData,
	}

	data, err := json.Marshal(errorMessage)

	if err == nil {
		logger.Debug("Sending error")
		sendStringTo(conn, string(data))
	} else {
		logger.Error("Error while sending error: " + err.Error())
	}
}

// sendMessageTo sends the data as normal message to the given connection
func sendMessageTo(connection *net.Conn, data string) error {
	message := Message{
		GenerallMessage: domain.GenerallMessage{
			MessageType: domain.MT_MESSAGE,
		},
		Data: data,
	}

	dataArray, err := json.Marshal(message)

	if err != nil {
		logger.Error("Error sending data: " + err.Error())
		return err
	}

	sendStringTo(connection, string(dataArray))
	return nil
}

// SendStringTo sends the given string with an \n character to the given connection.
func sendStringTo(connection *net.Conn, data string) {
	(*connection).Write([]byte(data + "\n"))
}
