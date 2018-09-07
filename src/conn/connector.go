package conn

import (
	"net"
	"sync"

	"github.com/go-messaging-service/goms-server/src/config"
	"github.com/go-messaging-service/goms-server/src/handler"
	"github.com/hauke96/sigolo"
)

type Connector struct {
	topics            []string
	connectionHandler []*handler.Handler
	initialized       bool
	mutex             *sync.Mutex
	distributor       *handler.Distributor
}

// Init will initialize the connection service by creating all topic notifier and initializing fields.
func (cs *Connector) Init(topics []string) {
	sigolo.Debug("Init connection service")

	cs.distributor = &handler.Distributor{}
	cs.distributor.Init(topics)

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
	cs.distributor.Add(&connHandler)
	cs.connectionHandler = append(cs.connectionHandler, &connHandler)
	cs.unlock()
	connHandler.HandleConnection()

	cs.lock()

	// TODO move whole list of handlers and this removal into distributor
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

// lock will prevent race conditions by ensuring that only one goroutine will have access to its fields.
func (cs *Connector) lock() {
	cs.mutex.Lock()
}

// unlock will free the fields so that other goroutines will have access to them.
func (cs *Connector) unlock() {
	cs.mutex.Unlock()
}
