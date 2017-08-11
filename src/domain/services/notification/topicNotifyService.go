package notificationServices

import (
	"encoding/json"
	"goms-server/src/domain/material"
	"goms-server/src/domain/services/common"
	technical "goms-server/src/technical/material"
	"goms-server/src/technical/services/logger"
	"sync"
)

type TopicNotifyService struct {
	Queue       chan *technical.Notification
	Errors      chan *technical.Notification
	Exit        chan bool
	initialized bool
	mutex       *sync.Mutex
}

// Init creates all neccessary channel (queues) to handle notifications.
func (tn *TopicNotifyService) Init() {
	tn.Queue = make(chan *technical.Notification)
	tn.Errors = make(chan *technical.Notification)
	tn.Exit = make(chan bool)

	tn.mutex = &sync.Mutex{}

	tn.initialized = true
}

// StartNotifier listens to incoming notification requests.
func (tn *TopicNotifyService) StartNotifier() {
	if !tn.initialized {
		logger.Fatal("TopicNotifyService not initialized")
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
func (tn *TopicNotifyService) sendNotification(notification *technical.Notification) {
	message := material.Message{
		Messagetype: material.MT_MESSAGE,
		Topics:      []string{notification.Topic},
		Data:        notification.Data,
	}

	if len(notification.Data) > 10 {
		logger.Info("send message with data: " + notification.Data[0:10] + "[...]")
	} else {
		logger.Info("send message with data: " + notification.Data)
	}

	messageByteArray, err := json.Marshal(message)
	messageString := string(messageByteArray)

	if err != nil {
		logger.Error("Error parsing message data: " + err.Error())
		for _, connection := range *notification.Connections {
			commonServices.SendErrorMessage(connection, material.ERR_SEND_FAILED, err.Error())
		}
		return
	}

	for _, connection := range *notification.Connections {
		//no error handling here, because we wouln't be able to send it to the client because SendError uses SendString
		commonServices.SendStringTo(connection, messageString)
	}
}
