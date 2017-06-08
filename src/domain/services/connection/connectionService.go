package connectionServices

import (
	"goMS/src/domain/material"
	"goMS/src/domain/services/common"
	"goMS/src/domain/services/notification"
	"goMS/src/technical/common"
	technical "goMS/src/technical/material"
	"goMS/src/technical/services/logger"
	"net"
	"sync"
)

const MAX_WAITING_CONNECTIONS = 1000

type ErrorMessage material.ErrorMessage

type ConnectionService struct {
	topics                      []string
	topicToConnection           map[string][]connectionHandler
	topicToNotificationServices map[string]notificationServices.TopicNotifyService
	initialized                 bool
	mutex                       *sync.Mutex
	ConnectionChannel           chan *net.Conn
}

// Init will initialize the connection service by creating all topic notifier and initializing fields.
func (cs *ConnectionService) Init(topics []string) {
	logger.Debug("Init connection service")

	cs.topicToConnection = make(map[string][]connectionHandler)

	cs.topicToNotificationServices = make(map[string]notificationServices.TopicNotifyService)
	for _, topic := range topics {
		service := notificationServices.TopicNotifyService{}
		service.Init()

		cs.topicToNotificationServices[topic] = service
		logger.Debug("Start notifier for " + topic)
		go func(service notificationServices.TopicNotifyService) {
			err := service.StartNotifier()

			if err != nil {
				logger.Fatal(err.Error()) // TODO really fatal here? Check better solutions (trying again? just print error?)
			}
		}(service)
	}

	cs.topics = topics
	cs.mutex = &sync.Mutex{}
	cs.ConnectionChannel = make(chan *net.Conn, MAX_WAITING_CONNECTIONS)

	cs.initialized = true
}

// Run listens to the port of this service and will start the handler.
func (cs *ConnectionService) Run() {
	if !cs.initialized {
		logger.Fatal("Connection Service not initialized!")
	}

	for {
		conn := <-cs.ConnectionChannel
		go cs.createAndRunHandler(conn)
	}
}

// createAndRunHandler sets up a new connection handler by registering to its events and starts it then.
// This should run on a new goroutine.
func (cs *ConnectionService) createAndRunHandler(conn *net.Conn) {
	logger.Debug("Create connection handler")

	connHandler := connectionHandler{
		connection: conn,
	}

	connHandler.RegisterEvent = append(connHandler.RegisterEvent, cs.handleRegisterEvent)
	connHandler.UnregisterEvent = append(connHandler.UnregisterEvent, cs.handleUnregisterEvent)
	connHandler.SendEvent = append(connHandler.SendEvent, cs.handleSendEvent)
	connHandler.HandleConnection()

	(*conn).Close()
}

// handleRegisterEvent should be called when a connection registered itself to a topic.
// This will return an error to the client when he wants to register to a topic he's not allowed to register him to.
func (cs *ConnectionService) handleRegisterEvent(conn connectionHandler, topics []string) {
	cs.lock()

	// A comma separated list of all topics, the client is not allowed to register to
	forbiddenTopics := ""
	alreadyRegisteredTopics := ""

	for _, topic := range topics {
		//TODO create a service for this. This should later take care of different user rights
		if !technicalCommon.ContainsString(cs.topics, topic) {
			forbiddenTopics += topic + ","
			logger.Info("Clients wants to register on invalid topic (" + topic + ").")

		} else if cs.isAlreadyRegistered(conn, topic) {
			alreadyRegisteredTopics += topic + ","
			logger.Debug("Client already registered on " + topic)

		} else {
			cs.topicToConnection[topic] = append(cs.topicToConnection[topic], conn)
			logger.Debug("Register " + topic)

		}
	}

	if len(forbiddenTopics) != 0 {
		commonServices.SendErrorMessage(conn.connection, material.ERR_REG_FORBIDDEN, forbiddenTopics)
	}
	if len(alreadyRegisteredTopics) != 0 {
		commonServices.SendErrorMessage(conn.connection, material.ERR_REG_ALREADY_REGISTERED, alreadyRegisteredTopics)
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

// remove will remove the given connection handler from the given array of handlers.
func (cs *ConnectionService) isAlreadyRegistered(h connectionHandler, topic string) bool {
	for _, a := range cs.topicToConnection[topic] {
		if a.connection == h.connection {
			return true
		}
	}
	return false
}
