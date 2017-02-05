package services_test

import (
	"bufio"
	"goMS/src/services"
	"goMS/src/technical/material"
	"goMS/src/technical/services/logger"
	"net"
	"testing"
	"time"
)

var conn1, conn2 *net.Conn
var buf1, buf2 *bufio.Reader
var serviceUnderTest *services.TopicNotifyService

func initConnections(t *testing.T) net.Listener {
	logger.TestMode = true

	//Create connections
	listener := listen(t)
	conn1, buf1 = dial(t)
	conn2, buf2 = dial(t)

	return listener
}

func tearDownConnection(l net.Listener) {
	l.Close()
	(*conn1).Close()
	(*conn2).Close()
}

func initNotifyService(t *testing.T) {
	serviceUnderTest = new(services.TopicNotifyService)
	serviceUnderTest.Init()
	go serviceUnderTest.StartNotifier()
}

func TestNotifyCorrectly(t *testing.T) {
	initNotifyService(t)
	l := initConnections(t)
	defer tearDownConnection(l)

	connections := make([]*net.Conn, 2)
	connections[0] = conn1
	connections[1] = conn2

	notification := technicalMaterial.Notification{
		Connections: &connections,
		Data:        "test123\n",
	}

	serviceUnderTest.Queue <- &notification

}

func TestNotInitializedCreatesError(t *testing.T) {
	serviceUnderTest = new(services.TopicNotifyService)
	// This is missing: serviceUnderTest.Init()
	// There must be an error here:
	err := serviceUnderTest.StartNotifier()

	if err == nil {
		t.Fatal("The service should return an error.")
	}
}

func TestSendToExitChanWillExitCorrectly(t *testing.T) {
	serviceUnderTest = new(services.TopicNotifyService)
	serviceUnderTest.Init()

	go func(service *services.TopicNotifyService, t *testing.T) {
		err := service.StartNotifier()

		if err != nil {
			t.Fatal()
		}
	}(serviceUnderTest, t)

	// Do we need this?
	time.Sleep(time.Millisecond)

	serviceUnderTest.Exit <- true
}

func listen(t *testing.T) net.Listener {
	l, err := net.Listen("tcp", ":3001")

	if err != nil {
		t.Fatal(err)
	}

	go func(l *net.Listener) {
		for {
			_, err := (*l).Accept()
			if err != nil {
				t.Fatal(err)
			}
		}
	}(&l)

	return l
}

func dial(t *testing.T) (*net.Conn, *bufio.Reader) {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(conn)

	return &conn, reader
}
