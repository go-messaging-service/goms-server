package services

import (
	"encoding/json"
	"goMS/src/material"
	"goMS/src/technical/common"
	"goMS/src/technical/services/logger"
	"net"
	"os"
	"strconv"
)

type ErrorMessage material.ErrorMessage

type ConnectionService struct {
	topics            []string
	topicToConnection map[string][]connectionHandler
	listener          net.Listener
	initialized       bool
}

func (cs *ConnectionService) Init(host string, port int, topics []string) {
	cs.topicToConnection = make(map[string][]connectionHandler)
	cs.listenTo(host, strconv.Itoa(port))

	cs.topics = topics

	cs.initialized = true
}

func (cs *ConnectionService) Run() {
	if !cs.initialized {
		logger.Error("Connection Service not initialized!")
		os.Exit(1)
	}

	for {
		conn, err := cs.waitForConnection()

		if err == nil {
			go cs.createAndRunHandler(conn)
		} else {
			logger.Error(err.Error())
		}
	}
}

func (cs *ConnectionService) createAndRunHandler(conn *net.Conn) {
	logger.Info("Create connection handler")

	connHandler := connectionHandler{
		connection: conn,
	}

	connHandler.RegisterEvent = append(connHandler.RegisterEvent, cs.handleRegisterEvent)
	connHandler.UnregisterEvent = append(connHandler.UnregisterEvent, cs.handleUnregisterEvent)
	connHandler.SendEvent = append(connHandler.SendEvent, cs.handleSendEvent)
	connHandler.HandleConnection()
}

func (cs *ConnectionService) handleRegisterEvent(conn connectionHandler, topics []string) {
	forbiddenTopics := ""

	for _, topic := range topics {
		if common.ContainsString(cs.topics, topic) {
			cs.topicToConnection[topic] = append(cs.topicToConnection[topic], conn)
			logger.Debug("Register " + topic)
		} else {
			//TODO send error message (or collect invalid topics to send one big message)
			forbiddenTopics += "," + topic
			logger.Info("Clients wants to register on invalid topic (" + topic + ").")
		}
	}

	if len(forbiddenTopics) != 0 {
		cs.sendErrorMessage(conn.connection, material.ERR_REG_FORBIDDEN, forbiddenTopics)
	}
}

func (cs *ConnectionService) handleUnregisterEvent(conn connectionHandler, topics []string) {
	for key, handlerList := range cs.topicToConnection {
		cs.topicToConnection[key] = remove(handlerList, conn)
	}
}

func (cs *ConnectionService) handleSendEvent(handler connectionHandler, topics []string, data string) {
	for _, topic := range topics {
		// Get all connections (as *net.Conn slice)
		handlerList := cs.topicToConnection[topic]
		connectionList := make([]*net.Conn, len(handlerList))
		for _, handler := range handlerList {
			connectionList = append(connectionList, handler.connection)
		}

		// create notification
		//		notification := &Notification{
		//			Connections: &connectionList,
		//			Data:        data,
		//		}

		//TODO send to notification services channel for messages
	}
}

func (cs *ConnectionService) sendErrorMessage(conn *net.Conn, errorCode, errorData string) {

	errorMessage := ErrorMessage{
		GenerallMessage: material.GenerallMessage{
			MessageType: material.MT_ERROR,
		},
		ErrorCode: errorCode,
		Error:     errorData,
	}

	data, err := json.Marshal(errorMessage)

	if err == nil {
		logger.Debug("Sending error")
		cs.sendStringTo(conn, string(data))
	} else {
		logger.Error("Error while sending error: " + err.Error())
	}
}

func (cs *ConnectionService) sendMessageTo(connection *net.Conn, data string) error {
	message := Message{
		GenerallMessage: material.GenerallMessage{
			MessageType: material.MT_MESSAGE,
		},
		Data: data,
	}

	dataArray, err := json.Marshal(message)

	if err != nil {
		logger.Error("Error sending data: " + err.Error())
		return err
	}

	cs.sendStringTo(connection, string(dataArray))
	return nil
}

func (cs *ConnectionService) sendStringTo(connection *net.Conn, data string) {
	(*connection).Write([]byte(data + "\n"))
}

func (cs *ConnectionService) listenTo(host, port string) {
	logger.Info("Try to listen on port " + port)

	listener, err := net.Listen("tcp", host+":"+port)

	if err == nil && listener != nil {
		logger.Info("Got listener for port " + port)
		cs.listener = listener
	} else if err != nil {
		logger.Error(err.Error())
		logger.Error("Maybe the port is not free?")
		os.Exit(1)
	} else if listener == nil {
		logger.Error("Could not listen to " + host + ":" + port + ". Unfortunately there's no error I could print here :( Check if no other services are running on port " + port + ".")
		os.Exit(1)
	}
}

func (cs *ConnectionService) waitForConnection() (*net.Conn, error) {
	conn, err := cs.listener.Accept()

	if err == nil {
		logger.Info("Got connection :D")
		return &conn, nil
	}

	logger.Error(err.Error())
	return nil, err
}

func remove(s []connectionHandler, e connectionHandler) []connectionHandler {
	for i, a := range s {
		if a.connection == e.connection {
			// Remove element at inedx i (s. "Slice Tricks" on github)
			// https://github.com/golang/go/wiki/SliceTricks
			logger.Debug("Remove element")
			s = append(s[:i], s[i+1:]...)
			return s
		}
	}
	return s
}
