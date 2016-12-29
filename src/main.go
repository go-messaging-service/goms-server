package main

import (
	"./logger"
	"net"
)

//TODO put this into kind of "connectionManager" or something
var topicToConnection map[string]net.Conn

func main() {
	logger.Info("Server starting...")

	topicToConnection = make(map[string]net.Conn)
	listenToPort("55545")

	logger.Info("Server startet")

	for {
		conn, err := waitForConnection()

		if err == nil {

			logger.Info("Create connection handler")

			connHandler := connectionHandler{
				connection: conn,
			}

			//TODO pass some lambdas from the connection manager to be able to send messages and receive them
			connHandler.HandleConnection()

		} else {
			logger.Err(err.Error())
		}
	}
}
