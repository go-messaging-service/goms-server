package conn

import (
	"net"
	"sync"

	"github.com/go-messaging-service/goms-server/src/config"
	"github.com/go-messaging-service/goms-server/src/dist"
	"github.com/go-messaging-service/goms-server/src/handler"
	"github.com/go-messaging-service/goms-server/src/msg"
	"github.com/hauke96/sigolo"
)

type Connector struct {
	topics                      []string
	connectionHandler           []*handler.Handler
	topicToNotificationServices map[string]dist.Notifier
	initialized                 bool
	mutex                       *sync.Mutex
}

// Init will initialize the connection service by creating all topic notifier and initializing fields.
func (cs *Connector) Init(topics []string) {
	sigolo.Debug("Init connection service")

	cs.topicToNotificationServices = make(map[string]dist.Notifier)
	for _, topic := range topics {
		service := dist.Notifier{}
		service.Init()

		cs.topicToNotificationServices[topic] = service

		sigolo.Debug("Start notifier for " + topic)

		go func(service dist.Notifier) {
			err := service.StartNotifier()

			if err != nil {
				sigolo.Fatal(err.Error())
			}
		}(service)
	}

	cs.topics = topics
	cs.mutex = &sync.Mutex{}

	cs.initialized = true
}

//HandleConnectionAsync creates a handler for the given connection and runs it in the background.
func (cs *Connector) HandleConnectionAsync(conn *net.Conn, config *config.Config) {
	go cs.createAndRunHandler(conn, config)
}

// createAndRunHandler sets up a new connection handler by registering to its events and starts it then.
// This should run on a new goroutine.
func (cs *Connector) createAndRunHandler(conn *net.Conn, config *config.Config) {
	sigolo.Debug("Create connection handler")

	connHandler := handler.Handler{}
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

// TODO move this according to https://github.com/go-messaging-service/goms-server/issues/2 into own service called "distributor" which is used by the handler. So the handler fires an event to the distributor and not to the connector.
// handleSendEvent sends the given data to all clients registeres to the given topics.
func (cs *Connector) handleSendEvent(handler handler.Handler, message *msg.Message) {
	//TODO move the lock into loop or is this a root for performance issues?
	cs.lock()
	for _, topic := range message.Topics {
		// Get all connections (as *net.Conn slice)
		var connectionList []*net.Conn

		for _, h := range cs.connectionHandler {
			if h.Connection != handler.Connection && h.IsRegisteredTo(topic) {
				connectionList = append(connectionList, h.Connection)
			}
		}

		// create notification
		notification := &dist.Notification{
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
func (cs *Connector) lock() {
	cs.mutex.Lock()
}

// unlock will free the fields so that other goroutines will have access to them.
func (cs *Connector) unlock() {
	cs.mutex.Unlock()
}
