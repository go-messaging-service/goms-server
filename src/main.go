package main

import (
	domainServices "goms-server/src/domain/services/connection"
	"goms-server/src/technical/material"
	"goms-server/src/technical/services"
	"goms-server/src/technical/services/logger"
	"net"
)

const VERSION string = "v0.3"

func main() {
	logger.Plain("           ,")
	logger.Plain("         ,/#/")
	logger.Plain("       ,/#/")
	logger.Plain("     ,/#/")
	logger.Plain("   ,/#/")
	logger.Plain(" ,/#/")
	logger.Plain("/#/___________________")
	logger.Plain("\\####################/")
	logger.Plain("  \\################/")
	logger.Plain("    \\############/")
	logger.Plain("      \\########/")
	logger.Plain("        \\####/")
	logger.Plain("          \\/")
	logger.Plain("")
	logger.Plain("Starting goMS " + VERSION + " ...")
	logger.Plain("I will just initialize myself and serve you as you configured me :)\n\n")

	logger.Info("Load configuration")
	config := loadConfig()

	logger.Info("Initialize logger")
	logger.DebugMode = config.ServerConfig.DebugLogging

	startServer(&config)
}

// startServer loads all configurations inits the services and starts them
func startServer(config *technicalMaterial.Config) {
	logger.Info("Initialize services")

	listeningServices := initConnectionService(config)

	logger.Info("Start connection listener")
	for _, listeningService := range listeningServices {
		go func(listeningService domainServices.ListeningService) {
			//TODO evaluate the need of a routine that restarts the service automatically when a error occurred. Something like: Error occurrec --> wait 5 seconds --> create service --> call Run()
			listeningService.Run()
		}(listeningService)
	}

	logger.Plain("\nThere we go, I'm ready to server ... eh ... serve\n")

	//TODO remove this and pass channels for closing
	select {}
}

// loadConfig loads the server config and its topics config.
func loadConfig() technicalMaterial.Config {
	logger.Info("Load configs")

	configLoader := technicalServices.ConfigLoader{}
	configLoader.LoadConfig("./conf/server.json")

	return configLoader.GetConfig()
}

// initConnectionService creates connection services bases on the given configuration.
func initConnectionService(config *technicalMaterial.Config) []domainServices.ListeningService {
	logger.Info("Initialize connection services")

	amountConnectors := len(config.ServerConfig.Connectors)

	listeningServices := make([]domainServices.ListeningService, amountConnectors)

	for i, connector := range config.ServerConfig.Connectors {
		// connection service
		connectionService := domainServices.ConnectionService{}
		connectionService.Init(config.TopicConfig.Topics)

		// listening service
		newConnectionClosure := func(conn *net.Conn) {
			connectionService.HandleConnectionAsync(conn, config)
		}
		listeningService := domainServices.ListeningService{}
		listeningService.Init(connector.Ip, connector.Port, config.TopicConfig.Topics, newConnectionClosure)

		listeningServices[i] = listeningService
	}

	return listeningServices
}
