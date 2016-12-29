package main

import (
	"./logger"
	"./services"
)

func main() {
	logger.Info("Server starting...")

	topicToConnection = make(map[string]net.Conn)
	listenToPort("55545")
	connectionManager := connectionManager{}

	logger.Info("Server startet")

	connectionManager.run()
}
