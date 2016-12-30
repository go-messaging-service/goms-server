package main

import (
	domainServices "goMS/src/services"
	technicalServices "goMS/src/technical/services"
	"goMS/src/technical/services/logger"
)

func main() {
	logger.DebugMode = true
	logger.Info("Server starting...")

	logger.Info("Load configs")
	configLoader := technicalServices.ConfigLoader{}
	configLoader.LoadTopics("./conf/topics.json")

	logger.Info("Initialize connection service")
	connectionService := domainServices.ConnectionService{}
	connectionService.Init("127.0.0.1", 55545, configLoader.TopicConfig.Topics)

	logger.Info("Server startet")

	connectionService.Run()
}
