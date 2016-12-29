package services

import (
	"../logger"
	"net"
	"os"
	"strconv"
)

type ConnectionService struct {
	topicToConnection map[string]net.Conn
	listener          net.Listener
	initialized       bool
}

func (cm *ConnectionService) Init(host string, port int) {
	cm.topicToConnection = make(map[string]net.Conn)
	cm.listenTo(host, strconv.Itoa(port))

	cm.initialized = true
}

func (cm *ConnectionService) Run() {
	if !cm.initialized {
		logger.Error("Connection Service not initialized!")
		os.Exit(1)
	}

	for {
		conn, err := cm.waitForConnection()

		if err == nil {

			logger.Info("Create connection handler")

			connHandler := connectionHandler{
				connection: conn,
			}

			connHandler.HandleConnection()

		} else {
			logger.Error(err.Error())
		}
	}
}

func (cm *ConnectionService) listenTo(host, port string) {
	logger.Info("Try to listen on port " + port)

	listener, err := net.Listen("tcp", host+":"+port)
	cm.listener = listener

	if err == nil {
		logger.Info("Got listener for port " + port)
	} else {
		logger.Error(err.Error())
	}
}

func (cm *ConnectionService) waitForConnection() (*net.Conn, error) {
	conn, err := cm.listener.Accept()

	if err == nil {
		logger.Info("Got connection :D")
		return &conn, nil
	}

	logger.Error(err.Error())
	return nil, err
}
