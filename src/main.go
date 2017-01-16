package main

import (
	domainServices "goMS/src/services"
	technicalServices "goMS/src/technical/services"
	"goMS/src/technical/services/logger"
)

func main() {
	logger.DebugMode = true
	logger.Plain("Welcome to the goMS (go Message Service)!")
	logger.Plain("I will just initialize me and serve you as you configured me :)\n\n")

	logger.Info("Initialize logger")

	logger.Info("Initialize server")
	configLoader := loadConfigs()
	connectionService := initConnectionService(configLoader.TopicConfig.Topics)

	logger.Info("Start server")
	connectionService.Run()
}

func loadConfigs() technicalServices.ConfigLoader {
	logger.Info("Load configs")

	configLoader := technicalServices.ConfigLoader{}
	configLoader.LoadTopics("./conf/topics.json")

	return configLoader
}

func initConnectionService(topics []string) domainServices.ConnectionService {
	logger.Info("Initialize connection service")

	connectionService := domainServices.ConnectionService{}
	connectionService.Init("127.0.0.1", 55545, topics)

	return connectionService
}
