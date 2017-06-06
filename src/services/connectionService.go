package services

import (
	"goMS/src/material"
	"goMS/src/technical/common"
	technical "goMS/src/technical/material"
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

// Init will initialize the connection service by creating all topic notifier and initializing fields.
func (cs *ConnectionService) Init(host string, port int, topics []string) {
	logger.Info("Init connection service for " + host + ":" + strconv.Itoa(port))

	cs.topicToConnection = make(map[string][]connectionHandler)

	cs.topicToNotificationServices = make(map[string]TopicNotifyService)
	for _, topic := range topics {
		service := TopicNotifyService{}
		service.Init()

		cs.topicToNotificationServices[topic] = service
		logger.Info("Start notifier for " + topic)
		go func(service TopicNotifyService) {
			err := service.StartNotifier()

			if err != nil {
				logger.Fatal(err.Error()) // TODO really fatal here? Check better solutions (trying again? just print error?)
			}
		}(service)
	}

	cs.topics = topics
	cs.host = host
	cs.port = strconv.Itoa(port)

	cs.mutex = &sync.Mutex{}

	cs.initialized = true
}

// Run listens to the port of this service and will start the handler.
func (cs *ConnectionService) Run() {
	if !cs.initialized {
		logger.Fatal("Connection Service not initialized!")
	}

	cs.listenTo(cs.host, cs.port)

	for {
		conn, err := cs.waitForConnection()

		if err == nil {
			go cs.createAndRunHandler(conn)
		} else {
			logger.Error(err.Error())
		}
	}
}

// createAndRunHandler sets up a new connection handler by registering to its events and starts it then.
// This should run on a new goroutine.
func (cs *ConnectionService) createAndRunHandler(conn *net.Conn) {
	logger.Info("Create connection handler")

	connHandler := connectionHandler{
		connection: conn,
	}

	connHandler.RegisterEvent = append(connHandler.RegisterEvent, cs.handleRegisterEvent)
	connHandler.UnregisterEvent = append(connHandler.UnregisterEvent, cs.handleUnregisterEvent)
	connHandler.SendEvent = append(connHandler.SendEvent, cs.handleSendEvent)
	connHandler.HandleConnection()

	(*conn).Close()
}

// listenTo actually listens to the port on the given host. It'll also exits the application if there's any problem.
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

// waitForConnection accepts an incoming connection request.
func (cs *ConnectionService) waitForConnection() (*net.Conn, error) {
	conn, err := cs.listener.Accept()

	if err == nil {
		logger.Info("Got connection :D")
		return &conn, nil
	}

	logger.Error(err.Error())
	return nil, err
}

// handleRegisterEvent should be called when a connection registered itself to a topic.
// This will return an error to the client when he wants to register to a topic he's not allowed to register him to.
func (cs *ConnectionService) handleRegisterEvent(conn connectionHandler, topics []string) {
	cs.lock()

	// A comma separated list of all topics, the client is not allowed to register to
	forbiddenTopics := ""

	for _, topic := range topics {
		//TODO create a service for this. This should later take care of different user rights
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

// handleUnregisterEvent unregisteres the client from the given topics. If there's a topic he's not registered to, nothing happens.
func (cs *ConnectionService) handleUnregisterEvent(conn connectionHandler, topics []string) {
	cs.lock()

	for topic, handlerList := range cs.topicToConnection {
		if technicalCommon.ContainsString(topics, topic) {
			cs.topicToConnection[topic] = remove(handlerList, conn)
		}
	}

	cs.unlock()
}

// handleSendEvent sends the given data to all clients registeres to the given topics.
func (cs *ConnectionService) handleSendEvent(handler connectionHandler, topics []string, data string) {
	for _, topic := range topics {
		// Get all connections (as *net.Conn slice)
		handlerList := cs.topicToConnection[topic]
		connectionList := make([]*net.Conn, len(handlerList))
		for i, handler := range handlerList {
			connectionList[i] = handler.connection
		}

		// create notification
		notification := &technical.Notification{
			Connections: &connectionList,
			Data:        data,
		}

		// puts the notification in the queue of the responsible service
		cs.topicToNotificationServices[topic].Queue <- notification
	}
}

// lock will prevent race conditions by ensuring that only one goroutine will have access to its fields.
func (cs *ConnectionService) lock() {
	cs.mutex.Lock()
}

// unlock will free the fields so that other goroutines will have access to them.
func (cs *ConnectionService) unlock() {
	cs.mutex.Unlock()
}

// remove will remove the given connection handler from the given array of handlers.
func remove(s []connectionHandler, e connectionHandler) []connectionHandler {
	result := []connectionHandler{}

	for _, a := range s {
		if a.connection != e.connection {
			result = append(result, a)
		}
	}

	return result
}
