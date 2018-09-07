package dist

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/go-messaging-service/goms-server/src/msg"
	technical "github.com/go-messaging-service/goms-server/src/technical/material"
	"github.com/hauke96/sigolo"
)

type Notifier struct {
	Queue       chan *technical.Notification
	Errors      chan *technical.Notification
	Exit        chan bool
	initialized bool
	mutex       *sync.Mutex
}

// Init creates all neccessary channel (queues) to handle notifications.
func (tn *Notifier) Init() {
	tn.Queue = make(chan *technical.Notification)
	tn.Errors = make(chan *technical.Notification)
	tn.Exit = make(chan bool)

	tn.mutex = &sync.Mutex{}

	tn.initialized = true
}

// StartNotifier listens to incoming notification requests.
func (tn *Notifier) StartNotifier() error {
	if !tn.initialized {
		return errors.New("TopicNotifyService not initialized")
	}

	for {
		select {
		case notification := <-tn.Queue:
			go tn.sendNotification(notification)
		case <-tn.Exit:
			break
		}
	}
}

// sendNotification sends the notification or an error if there's one.
func (tn *Notifier) sendNotification(notification *technical.Notification) {
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
