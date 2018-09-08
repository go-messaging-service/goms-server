package dist

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/go-messaging-service/goms-server/src/msg"
	"github.com/hauke96/sigolo"
)

// SendStringTo sends the given string with an \n character to the given connection.
func SendStringTo(connection *net.Conn, data string) error {
	_, err := (*connection).Write([]byte(data + "\n"))

	return err
}

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

// SendMessage enqueues the request to send the message into the queue.
// Therefore the sending itself will happen a bit later because the background
// thread will read from the queue.
func (n *Notifier) SendMessage(connections []*net.Conn, topic, message string) {
	// create notification
	notification := &Notification{
		Connections: &connections,
		Topic:       topic,
		Data:        message,
	}

	n.Queue <- notification
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
			n.SendError(connection, msg.ERR_SEND_FAILED, err.Error())
		}
		return
	}

	for _, connection := range *notification.Connections {
		err := SendStringTo(connection, messageString)

		if err != nil {
			sigolo.Error(fmt.Sprintf("Could not send message to %s", (*connection).RemoteAddr()))
		}
	}
}

// TODO also create a queue for the errors
// SendError directly sends the error, there's no asynchronous queue here.
func (n *Notifier) SendError(connection *net.Conn, errorCode, message string) {
	errorMessage := msg.ErrorMessage{
		Messagetype: msg.MT_ERROR,
		Errorcode:   errorCode,
		Error:       message,
	}

	data, err := json.Marshal(errorMessage)

	if err == nil {
		sigolo.Debug("Sending error")
		SendStringTo(connection, string(data)+"\n")
	} else {
		sigolo.Error("Error while sending error: " + err.Error())
	}

}
