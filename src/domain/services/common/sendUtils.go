package commonServices

import (
	"encoding/json"
	"goMS/src/domain/material"
	"goMS/src/technical/services/logger"
	"net"
)

type ErrorMessage material.ErrorMessage
type Message material.Message

// sendErrorMessage sends the given error data as an error message to the given client.
func SendErrorMessage(conn *net.Conn, errorCode, errorData string) {

	errorMessage := ErrorMessage{
		AbstractMessage: material.AbstractMessage{
			MessageType: material.MT_ERROR,
		},
		ErrorCode: errorCode,
		Error:     errorData,
	}

	data, err := json.Marshal(errorMessage)

	if err == nil {
		logger.Debug("Sending error")
		SendStringTo(conn, string(data))
	} else {
		logger.Error("Error while sending error: " + err.Error())
	}
}

// sendMessageTo sends the data as normal message to the given connection
func SendMessageTo(connection *net.Conn, data string) error {
	message := Message{
		AbstractMessage: material.AbstractMessage{
			MessageType: material.MT_MESSAGE,
		},
		Data: data,
	}

	dataArray, err := json.Marshal(message)

	if err != nil {
		logger.Error("Error sending data: " + err.Error())
		return err
	}

	SendStringTo(connection, string(dataArray))
	return nil
}

// SendStringTo sends the given string with an \n character to the given connection.
func SendStringTo(connection *net.Conn, data string) {
	(*connection).Write([]byte(data + "\n"))
}
