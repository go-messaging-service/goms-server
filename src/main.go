package main

import (
	domainServices "goMS/src/services"
	technicalMaterial "goMS/src/technical/material"
	technicalServices "goMS/src/technical/services"
	"goMS/src/technical/services/logger"
)

func main() {
	logger.Info("Initialize logger")

	logger.DebugMode = true
	logger.Plain("Welcome to the goMS (go Message Service)!")
	logger.Plain("I will just initialize me and serve you as you configured me :)\n\n")

	startServer()
}

func startServer() {
	logger.Info("Initialize server")
	config := loadConfig()
	connectionServices := initConnectionService(config)

	logger.Info("Start server")
	for _, connectionService := range connectionServices {
		connectionService.Run()
	}

	//TODO remove this and pass channels for closing
	for {
	}
}

func loadConfig() technicalMaterial.Config {
	logger.Info("Load configs")

	configLoader := technicalServices.ConfigLoader{}
	configLoader.LoadConfig("./conf/server.json")

	return configLoader.GetConfig()
}

func initConnectionService(config technicalMaterial.Config) []domainServices.ConnectionService {
	logger.Info("Initialize connection service")

	connectionServices := make([]domainServices.ConnectionService, len(config.ServerConfig.Connectors))

	for i, connector := range config.ServerConfig.Connectors {
		connectionService := domainServices.ConnectionService{}
		connectionService.Init(connector.Ip, connector.Port, config.TopicConfig.Topics)

		connectionServices[i] = connectionService
	}

	return connectionServices
}
