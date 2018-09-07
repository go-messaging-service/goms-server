package commonServices

import (
	"encoding/json"
	"net"

	"github.com/go-messaging-service/goms-server/src/domain/material"
	"github.com/hauke96/sigolo"
)

// sendErrorMessage sends the given error data as an error message to the given client.
func SendErrorMessage(conn *net.Conn, errorCode, errorData string) {

	errorMessage := material.ErrorMessage{
		Messagetype: material.MT_ERROR,
		Errorcode:   errorCode,
		Error:       errorData,
	}

	data, err := json.Marshal(errorMessage)

	if err == nil {
		sigolo.Debug("Sending error")
		SendStringTo(conn, string(data))
	} else {
		sigolo.Error("Error while sending error: " + err.Error())
	}
}

// SendStringTo sends the given string with an \n character to the given connection.
func SendStringTo(connection *net.Conn, data string) {
	(*connection).Write([]byte(data + "\n"))
}
