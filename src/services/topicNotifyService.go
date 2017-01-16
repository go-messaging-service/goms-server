package services

import (
	"encoding/json"
	domain "goMS/src/material"
	technical "goMS/src/technical/material"
	"goMS/src/technical/services/logger"
	"net"
)

type Notification technical.Notification

type TopicNotifyService struct {
	queue  chan *Notification
	errors chan *Notification
	exit   chan bool
}

func (tn *TopicNotifyService) Init() {
	tn.queue = make(chan *Notification)
	tn.errors = make(chan *Notification)
	tn.exit = make(chan bool)
}

func (tn *TopicNotifyService) StartNotifier() {
	for {
		select {
		case notification := <-tn.queue:
			tn.sendNotification(notification)
		case <-tn.exit:
			break
		}
	}
}

func (tn *TopicNotifyService) sendNotification(notification *Notification) {
	//TODO implement sending to all connections

	//		err := cs.sendMessageTo(destHandler.connection, data)

	//		if err != nil {
	//			cs.sendErrorMessage(handler.connection, material.ERR_SEND_FAILED, err.Error())
	//		}
}

func (tn *TopicNotifyService) sendMessageTo(connection *net.Conn, data string) error {
	message := Message{
		GenerallMessage: domain.GenerallMessage{
			MessageType: domain.MT_MESSAGE,
		},
		Data: data,
	}

	dataArray, err := json.Marshal(message)

	if err != nil {
		logger.Error("Error sending data: " + err.Error())
		return err
	}

	tn.sendStringTo(connection, string(dataArray))
	return nil
}

func (tn *TopicNotifyService) sendStringTo(connection *net.Conn, data string) {
	(*connection).Write([]byte(data + "\n"))
}
