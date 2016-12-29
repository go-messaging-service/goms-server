package main

import (
	"./logger"
	"./services"
)

func main() {
	logger.Info("Server starting...")

	connectionManager := services.ConnectionManager{}
	connectionManager.Init("127.0.0.1", 55545)

	logger.Info("Server startet")

	connectionManager.Run()
}
