package notificationServices

import (
	domain "goMS/src/domain/material"
	"goMS/src/domain/services/common"
	technical "goMS/src/technical/material"
	"goMS/src/technical/services/logger"
)

type TopicNotifyService struct {
	Queue       chan *technical.Notification
	Errors      chan *technical.Notification
	Exit        chan bool
	initialized bool
}

// Init creates all neccessary channel (queues) to handle notifications.
func (tn *TopicNotifyService) Init() {
	tn.Queue = make(chan *technical.Notification)
	tn.Errors = make(chan *technical.Notification)
	tn.Exit = make(chan bool)

	tn.initialized = true
}

// StartNotifier listens to incoming notification requests.
func (tn *TopicNotifyService) StartNotifier() {
	// Not initialized
	if !tn.initialized {
		logger.Fatal("TopicNotifyService not initialized")
	}

	for {
		select {
		case notification := <-tn.Queue:
			tn.sendNotification(notification)
		case <-tn.Exit:
			break
		}
	}
}

// sendNotification sends the notification or an error if there's one.
func (tn *TopicNotifyService) sendNotification(notification *technical.Notification) {
	for _, connection := range *notification.Connections {

		err := commonServices.SendMessageTo(connection, notification.Data)

		if err != nil {
			commonServices.SendErrorMessage(connection, domain.ERR_SEND_FAILED, err.Error())
		}
	}
}
