package services

import (
	"../logger"
	"net"
)

type connectionManager struct {
	//TODO put this into kind of "connectionManager" or something
	topicToConnection map[string]net.Conn
}

func (cm *connectionManager) Run() {
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
