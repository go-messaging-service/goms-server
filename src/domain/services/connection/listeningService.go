package connectionServices

import (
	"goMS/src/technical/services/logger"
	"net"
	"strconv"
)

type ListeningService struct {
	connectionService ConnectionService
	listener          net.Listener
	initialized       bool
	host              string
	port              string
	connectionChannel chan *net.Conn
}

func (ls *ListeningService) Init(host string, port int, topics []string, connectionChannel chan *net.Conn) {
	logger.Debug("Init listening service for " + host + ":" + strconv.Itoa(port))

	ls.host = host
	ls.port = strconv.Itoa(port)
	ls.connectionChannel = connectionChannel

	ls.listenTo()

	ls.initialized = true
}

// Run listens to the port of this service and will start the handler.
func (ls *ListeningService) Run() {
	if !ls.initialized {
		logger.Fatal("Listening Service not initialized!")
	}

	for {
		conn, err := ls.waitForConnection()

		if err == nil {
			ls.connectionChannel <- conn
		} else {
			logger.Error(err.Error())
		}
	}
}

// listenTo actually listens to the port on the given host. It'll also exits the application if there's any problem.
func (ls *ListeningService) listenTo() {
	logger.Debug("Try to listen on port " + ls.port)

	listener, err := net.Listen("tcp", ls.host+":"+ls.port)

	if err == nil && listener != nil {
		logger.Debug("Got listener for port " + ls.port)
		ls.listener = listener
	} else if err != nil {
		logger.Error(err.Error())
		logger.Fatal("Maybe the port is not free?")
	} else if listener == nil {
		logger.Fatal("Could not listen to " + ls.host + ":" + ls.port + ". Unfortunately there's no error I could print here :( Check if no other services are running on port " + ls.port + ".")
	}
}

// waitForConnection accepts an incoming connection request.
func (ls *ListeningService) waitForConnection() (*net.Conn, error) {
	conn, err := ls.listener.Accept()

	if err == nil {
		logger.Info("Got connection on " + ls.host + ":" + ls.port)
		return &conn, nil
	}

	logger.Error(err.Error())
	return nil, err
}
