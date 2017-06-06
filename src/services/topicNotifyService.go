package services

import (
	"errors"
	domain "goMS/src/material"
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
func (tn *TopicNotifyService) StartNotifier() error {
	// Not initialized
	if !tn.initialized {
		logger.Error("TopicNotifyService not initialized")
		return errors.New("TopicNotifyService not initialized")
	}

	for {
		select {
		case notification := <-tn.Queue:
			tn.sendNotification(notification)
		case <-tn.Exit: //TODO make the value "true" that should be sent in this channel as global constant
			break
		}
	}

	return nil
}

// sendNotification sends the notification or an error if there's one.
func (tn *TopicNotifyService) sendNotification(notification *technical.Notification) {
	for _, connection := range *notification.Connections {

		err := sendMessageTo(connection, notification.Data)

		if err != nil {
			sendErrorMessage(connection, domain.ERR_SEND_FAILED, err.Error())
		}
	}
}
