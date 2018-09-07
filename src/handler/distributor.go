package handler

import (
	"net"
	"sync"

	"github.com/go-messaging-service/goms-server/src/dist"
	"github.com/go-messaging-service/goms-server/src/msg"
	"github.com/hauke96/sigolo"
)

type Distributor struct {
	knownHandler                []*Handler
	topicToNotificationServices map[string]dist.Notifier
	mutex                       *sync.Mutex
}

func (d *Distributor) Init(topics []string) {
	d.topicToNotificationServices = make(map[string]dist.Notifier)
	d.mutex = &sync.Mutex{}

	for _, topic := range topics {
		service := dist.Notifier{}
		service.Init()

		d.topicToNotificationServices[topic] = service

		sigolo.Debug("Start notifier for " + topic)

		go func(service dist.Notifier) {
			err := service.StartNotifier()

			if err != nil {
				sigolo.Fatal(err.Error())
			}
		}(service)
	}
}

func (cs *Distributor) HandleSendEvent(handler Handler, message *msg.Message) {
	//TODO move the lock into loop or is this a root for performance issues?
	cs.lock()
	for _, topic := range message.Topics {
		// Get all connections (as *net.Conn slice)
		var connectionList []*net.Conn

		for _, h := range cs.knownHandler {
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
		cs.topicToNotificationServices[topic].Queue <- notification
	}
	cs.unlock()
}

// lock will prevent race conditions by ensuring that only one goroutine will have access to its fields.
func (cs *Distributor) lock() {
	cs.mutex.Lock()
}

// unlock will free the fields so that other goroutines will have access to them.
func (cs *Distributor) unlock() {
	cs.mutex.Unlock()
}
