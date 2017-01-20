package services

import (
	"goMS/src/material"
	"goMS/src/technical/common"
	"goMS/src/technical/services/logger"
	"net"
	"strconv"
	"sync"
)

type ErrorMessage material.ErrorMessage

type ConnectionService struct {
	topics                      []string
	topicToConnection           map[string][]connectionHandler
	topicToNotificationServices map[string]TopicNotifyService
	listener                    net.Listener
	initialized                 bool
	host                        string
	port                        string
	mutex                       *sync.Mutex
}

func (cs *ConnectionService) Init(host string, port int, topics []string) {
	cs.topicToConnection = make(map[string][]connectionHandler)

	cs.topicToNotificationServices = make(map[string]TopicNotifyService)
	for _, topic := range topics {
		service := TopicNotifyService{}
		service.Init()

		cs.topicToNotificationServices[topic] = service
		logger.Info("Start notifier for " + topic)
		go service.StartNotifier()
	}

	cs.topics = topics
	cs.host = host
	cs.port = strconv.Itoa(port)

	cs.mutex = &sync.Mutex{}

	cs.initialized = true
}

func (cs *ConnectionService) Run() {
	if !cs.initialized {
		logger.Fatal("Connection Service not initialized!")
	}

	cs.listenTo(cs.host, cs.port)

	go func() {
		for {
			conn, err := cs.waitForConnection()

			if err == nil {
				go cs.createAndRunHandler(conn)
			} else {
				logger.Error(err.Error())
			}
		}
	}()
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

func (cs *ConnectionService) listenTo(host, port string) {
	logger.Info("Try to listen on port " + port)

	listener, err := net.Listen("tcp", host+":"+port)

	if err == nil && listener != nil {
		logger.Info("Got listener for port " + port)
		cs.listener = listener
	} else if err != nil {
		logger.Error(err.Error())
		logger.Fatal("Maybe the port is not free?")
	} else if listener == nil {
		logger.Fatal("Could not listen to " + host + ":" + port + ". Unfortunately there's no error I could print here :( Check if no other services are running on port " + port + ".")
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

func (cs *ConnectionService) handleRegisterEvent(conn connectionHandler, topics []string) {
	cs.lock()

	forbiddenTopics := ""

	for _, topic := range topics {
		if technicalCommon.ContainsString(cs.topics, topic) {
			cs.topicToConnection[topic] = append(cs.topicToConnection[topic], conn)
			logger.Debug("Register " + topic)
		} else {
			forbiddenTopics += "," + topic
			logger.Info("Clients wants to register on invalid topic (" + topic + ").")
		}
	}

	if len(forbiddenTopics) != 0 {
		sendErrorMessage(conn.connection, material.ERR_REG_FORBIDDEN, forbiddenTopics)
	}

	cs.unlock()
}

func (cs *ConnectionService) handleUnregisterEvent(conn connectionHandler, topics []string) {
	cs.lock()

	for key, handlerList := range cs.topicToConnection {
		cs.topicToConnection[key] = remove(handlerList, conn)
	}

	cs.unlock()
}

func (cs *ConnectionService) handleSendEvent(handler connectionHandler, topics []string, data string) {
	for _, topic := range topics {
		// Get all connections (as *net.Conn slice)
		handlerList := cs.topicToConnection[topic]
		connectionList := make([]*net.Conn, len(handlerList))
		for i, handler := range handlerList {
			connectionList[i] = handler.connection
		}

		// create notification
		notification := &Notification{
			Connections: &connectionList,
			Data:        data,
		}

		cs.topicToNotificationServices[topic].queue <- notification
	}
}

func (cs *ConnectionService) lock() {
	cs.mutex.Lock()
}

func (cs *ConnectionService) unlock() {
	cs.mutex.Unlock()
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
