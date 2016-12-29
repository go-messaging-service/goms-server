package main

import (
	"./logger"
	"./services"
)

func main() {
	logger.DebugMode = true
	logger.Info("Server starting...")

	connectionService := services.ConnectionService{}
	connectionService.Init("127.0.0.1", 55545)

	logger.Info("Server startet")

	connectionService.Run()
}
