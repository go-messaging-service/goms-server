package connectionServices

import (
	"goms-server/src/domain/material"
	"goms-server/src/domain/services/notification"
	technical "goms-server/src/technical/material"
	"goms-server/src/technical/services/logger"
	"net"
	"sync"
)

type ErrorMessage material.ErrorMessage

type ConnectionService struct {
	topics                      []string
	connectionHandler           []*connectionHandler
	topicToNotificationServices map[string]notificationServices.TopicNotifyService
	initialized                 bool
	mutex                       *sync.Mutex
}

// Init will initialize the connection service by creating all topic notifier and initializing fields.
func (cs *ConnectionService) Init(topics []string) {
	logger.Debug("Init connection service")

	cs.topicToNotificationServices = make(map[string]notificationServices.TopicNotifyService)
	for _, topic := range topics {
		service := notificationServices.TopicNotifyService{}
		service.Init()

		cs.topicToNotificationServices[topic] = service

		logger.Debug("Start notifier for " + topic)

		go func(service notificationServices.TopicNotifyService) {
			err := service.StartNotifier()

			if err != nil {
				logger.Fatal(err.Error())
			}
		}(service)
	}

	cs.topics = topics
	cs.mutex = &sync.Mutex{}

	cs.initialized = true
}

//HandleConnectionAsync creates a handler for the given connection and runs it in the background.
func (cs *ConnectionService) HandleConnectionAsync(conn *net.Conn, config *technical.Config){
	go cs.createAndRunHandler(conn, config)
}

// createAndRunHandler sets up a new connection handler by registering to its events and starts it then.
// This should run on a new goroutine.
func (cs *ConnectionService) createAndRunHandler(conn *net.Conn, config *technical.Config) {
	logger.Debug("Create connection handler")

	connHandler := connectionHandler{}
	connHandler.Init(conn, config)

	cs.lock()
	connHandler.SendEvent = append(connHandler.SendEvent, cs.handleSendEvent)
	cs.connectionHandler = append(cs.connectionHandler, &connHandler)
	cs.unlock()
	connHandler.HandleConnection()

	cs.lock()

	// find connection handler index
	i := -1
	for j, a := range cs.connectionHandler {
		if a == &connHandler {
			i = j
			break
		}
	}

	// remove connection handler
	if i != -1 {
		cs.connectionHandler = append(cs.connectionHandler[:i], cs.connectionHandler[i+1:]...)
	}

	cs.unlock()

	(*conn).Close()
}

// handleSendEvent sends the given data to all clients registeres to the given topics.
func (cs *ConnectionService) handleSendEvent(handler connectionHandler, message *Message) {
	//TODO move the lock into loop or is this a root for performance issues?
	cs.lock()
	for _, topic := range message.Topics {
		// Get all connections (as *net.Conn slice)
		var connectionList []*net.Conn

		for _, h := range cs.connectionHandler {
			if h.connection != handler.connection && h.isRegisteredTo(topic) {
				connectionList = append(connectionList, h.connection)
			}
		}

		// create notification
		notification := &technical.Notification{
			Connections: &connectionList,
			Topic:       topic,
			Data:        message.Data,
		}

		// puts the notification in the queue of the responsible service
		cs.topicToNotificationServices[topic].Queue <- notification
	}
	cs.unlock()
}

// lock will prevent race conditions by ensuring that only one goroutine will have access to its fields.
func (cs *ConnectionService) lock() {
	cs.mutex.Lock()
}

// unlock will free the fields so that other goroutines will have access to them.
func (cs *ConnectionService) unlock() {
	cs.mutex.Unlock()
}
