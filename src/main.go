package main

import (
	"goMS/src/services"
	"goMS/src/technical/services/logger"
)

func main() {
	logger.DebugMode = true
	logger.Info("Server starting...")

	logger.Info("Load configs")

	logger.Info("Initialize connection service")
	connectionService := services.ConnectionService{}
	connectionService.Init("127.0.0.1", 55545)

	logger.Info("Server startet")

	connectionService.Run()
}
