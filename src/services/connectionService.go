package services

import (
	"../logger"
	"net"
	"os"
	"strconv"
)

type ConnectionService struct {
	topicToConnection map[string][]connectionHandler
	listener          net.Listener
	initialized       bool
}

func (cs *ConnectionService) Init(host string, port int) {
	cm.topicToConnection = make(map[string][]connectionHandler)
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

			connHandler.RegisterEvent = append(connHandler.RegisterEvent, cm.handleRegisterEvent)
			connHandler.HandleConnection()

		} else {
			logger.Error(err.Error())
		}
	}
}

func (cm *ConnectionService) handleRegisterEvent(conn connectionHandler, topics []string) {
	for _, topic := range topics {
		cm.topicToConnection[topic] = append(cm.topicToConnection[topic], conn)
		logger.Debug("Register " + topic)
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
