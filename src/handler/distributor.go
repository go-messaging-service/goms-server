package handler

import (
	"net"
	"sync"

	"github.com/go-messaging-service/goms-server/src/dist"
	"github.com/go-messaging-service/goms-server/src/msg"
	"github.com/hauke96/sigolo"
)

type Distributor struct {
	knownHandler []*Handler
	notifier     *dist.Notifier
	mutex        *sync.Mutex
}

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

func (d *Distributor) Add(handler *Handler) {
	handler.SendEvent = append(handler.SendEvent, d.HandleSendEvent)
	handler.ErrorEvent = append(handler.ErrorEvent, d.HandleErrorEvent)
}

func (d *Distributor) HandleSendEvent(handler Handler, message *msg.Message) {
	//TODO move the lock into loop or is this a root for performance issues?
	d.lock()
	for _, topic := range message.Topics {
		// Get all connections (as *net.Conn slice)
		var connectionList []*net.Conn

		for _, h := range d.knownHandler {
			if h.connection != handler.connection && h.IsRegisteredTo(topic) {
				connectionList = append(connectionList, h.connection)
			}
		}

		// create notification
		notification := &dist.Notification{
			Connections: &connectionList,
			Topic:       topic,
			Data:        message.Data,
		}

		// puts the notification in the queue of the responsible service
		d.notifier.Queue <- notification
	}
	d.unlock()
}

// TODO maybe just pass connection instead of whole handler?
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
