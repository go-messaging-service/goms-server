package main

import (
	"net"
	"os"

	"github.com/go-messaging-service/goms-server/src/config"
	"github.com/go-messaging-service/goms-server/src/conn"
	"github.com/hauke96/kingpin"
	"github.com/hauke96/sigolo"
)

const VERSION string = "v0.4.0"

var (
	app           = kingpin.New("goMS", "A simple messaging service written in go")
	appConfigFile = app.Flag("config", "Specifies the configuration file that should be used. This is \"./conf/server.json\" by default.").Short('c').Default("./conf/server.json").String()
)

func main() {
	configureCLI()
	app.Parse(os.Args[1:])

	printWelcomeScreen()

	sigolo.Info("Load configuration")
	config := loadConfig()

	sigolo.Info("Initialize logger")
	if config.ServerConfig.DebugLogging {
		sigolo.LogLevel = sigolo.LOG_DEBUG
	} else {
		sigolo.LogLevel = sigolo.LOG_INFO
	}

	startServer(&config)
}

func configureCLI() {
	app.Author("Hauke Stieler")
	app.Version(VERSION)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')
}

func printWelcomeScreen() {
	sigolo.LogLevel = sigolo.LOG_PLAIN

	sigolo.Plain("           ,")
	sigolo.Plain("         ,/#/")
	sigolo.Plain("       ,/#/")
	sigolo.Plain("     ,/#/")
	sigolo.Plain("   ,/#/")
	sigolo.Plain(" ,/#/")
	sigolo.Plain("/#/___________________")
	sigolo.Plain("\\####################/")
	sigolo.Plain("  \\################/")
	sigolo.Plain("    \\############/")
	sigolo.Plain("      \\########/")
	sigolo.Plain("        \\####/")
	sigolo.Plain("          \\/")
	sigolo.Plain("")
	sigolo.Plain("Starting goMS " + VERSION + " ...")
	sigolo.Plain("I will just initialize myself and serve you as you configured me :)\n\n")
}

// startServer loads all configurations inits the services and starts them
func startServer(config *config.Config) {
	sigolo.Info("Initialize services")

	listeningServices := initConnectionService(config)

	sigolo.Info("Start connection listener")
	for _, listeningService := range listeningServices {
		go func(listeningService conn.Listener) {
			//TODO evaluate the need of a routine that restarts the service automatically when a error occurred. Something like: Error occurrec --> wait 5 seconds --> create service --> call Run()
			listeningService.Run()
		}(listeningService)
	}

	sigolo.Debug("\nThere we go, I'm ready to server ... eh ... serve\n")

	//TODO remove this and pass channels for closing
	select {}
}

// loadConfig loads the server config and its topics config.
func loadConfig() config.Config {
	sigolo.Info("Load configs")

	configLoader := config.ConfigLoader{}
	configLoader.LoadConfig(*appConfigFile)

	return configLoader.GetConfig()
}

// initConnectionService creates connection services bases on the given configuration.
func initConnectionService(config *config.Config) []conn.Listener {
	sigolo.Info("Initialize connection services")

	amountConnectors := len(config.ServerConfig.Connectors)

	listeningServices := make([]conn.Listener, amountConnectors)

	for i, connector := range config.ServerConfig.Connectors {
		// connection service
		connectionService := conn.Connector{}
		connectionService.Init(config.TopicConfig.Topics)

		// listening service
		newConnectionClosure := func(conn *net.Conn) {
			connectionService.HandleConnectionAsync(conn, config)
		}
		listeningService := conn.Listener{}
		listeningService.Init(connector.Ip, connector.Port, config.TopicConfig.Topics, newConnectionClosure)

		listeningServices[i] = listeningService
	}

	return listeningServices
}
