package dist

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/go-messaging-service/goms-server/src/msg"
	"github.com/hauke96/sigolo"
)

type Notifier struct {
	Queue       chan *Notification
	Errors      chan *Notification
	Exit        chan bool
	initialized bool
	mutex       *sync.Mutex
}

// Init creates all neccessary channel (queues) to handle notifications.
func (n *Notifier) Init() {
	n.Queue = make(chan *Notification)
	n.Errors = make(chan *Notification)
	n.Exit = make(chan bool)

	n.mutex = &sync.Mutex{}

	n.initialized = true
}

// StartNotifier listens to incoming notification requests.
func (n *Notifier) StartNotifier() error {
	if !n.initialized {
		return errors.New("TopicNotifyService not initialized")
	}

	for {
		select {
		case notification := <-n.Queue:
			go n.sendNotification(notification)
		case <-n.Exit:
			break
		}
	}
}

// sendNotification sends the notification or an error if there's one.
func (n *Notifier) sendNotification(notification *Notification) {
	message := msg.Message{
		Messagetype: msg.MT_MESSAGE,
		Topics:      []string{notification.Topic},
		Data:        notification.Data,
	}

	if len(notification.Data) > 10 {
		sigolo.Info("send message with data: " + notification.Data[0:10] + "[...]")
	} else {
		sigolo.Info("send message with data: " + notification.Data)
	}

	messageByteArray, err := json.Marshal(message)
	messageString := string(messageByteArray)

	if err != nil {
		sigolo.Error("Error parsing message data: " + err.Error())
		for _, connection := range *notification.Connections {
			SendErrorMessage(connection, msg.ERR_SEND_FAILED, err.Error())
		}
		return
	}

	for _, connection := range *notification.Connections {
		//no error handling here, because we wouln't be able to send it to the client because SendError uses SendString
		SendStringTo(connection, messageString)
	}
}
