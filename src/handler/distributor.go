package handler

import (
	"fmt"
	"net"
	"sync"

	"github.com/go-messaging-service/goms-server/src/dist"
	"github.com/hauke96/sigolo"
)

type Distributor struct {
	knownHandler []*Handler
	notifier     *dist.Notifier
	mutex        *sync.Mutex
}

// Init initializes the distributor and starts the notifier.
func (d *Distributor) Init(topics []string) {
	d.notifier = &dist.Notifier{}
	d.notifier.Init()
	d.mutex = &sync.Mutex{}

	sigolo.Debug("Start notifier")

	go func(service *dist.Notifier) {
		err := service.StartNotifier()

		if err != nil {
			sigolo.Fatal(err.Error())
		}
	}(d.notifier)
}

// Add adds the handler to the list of handlers and registeres the distributor to events.
func (d *Distributor) Add(handler *Handler) {
	d.knownHandler = append(d.knownHandler, handler)

	handler.SendEvent = append(handler.SendEvent, d.HandleSendEvent)
	handler.ErrorEvent = append(handler.ErrorEvent, d.HandleErrorEvent)
}

// HandleSendEvent will determine the receivers of the message and passes the
// request to send the message to the notifier.
func (d *Distributor) HandleSendEvent(handler Handler, topics []string, message string) {
	//TODO move the lock into loop or is this a root for performance issues?
	d.lock()
	for _, topic := range topics {
		// Get all connections (as *net.Conn slice)
		var connectionList []*net.Conn

		for _, h := range d.knownHandler {
			sigolo.Debug(fmt.Sprintf("Check handler %v for topic %s", h, topic))
			if h.connection != handler.connection && h.IsRegisteredTo(topic) {
				sigolo.Debug(fmt.Sprintf("Found handler %v for topic %s", h, topic))
				connectionList = append(connectionList, h.connection)
			}
		}

		d.notifier.SendMessage(connectionList, topic, message)
	}
	d.unlock()
}

// TODO maybe just pass connection instead of whole handler?
// HandleErrorEvent will pass the error to the notifier, who will send the error.
func (d *Distributor) HandleErrorEvent(handler *Handler, errorCode, message string) {
	d.notifier.SendError(handler.connection, errorCode, message)
}

// lock will prevent race conditions by ensuring that only one goroutine will have access to its fields.
func (d *Distributor) lock() {
	d.mutex.Lock()
}

// unlock will free the fields so that other goroutines will have access to them.
func (d *Distributor) unlock() {
	d.mutex.Unlock()
}
