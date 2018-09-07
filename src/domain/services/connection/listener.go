package connectionServices

import (
	"net"
	"strconv"

	"github.com/hauke96/sigolo"
)

type Listener struct {
	listener          net.Listener
	initialized       bool
	host              string
	port              string
	connectionChannel func(*net.Conn)
}

func (ls *Listener) Init(host string, port int, topics []string, connectionChannel func(*net.Conn)) {
	sigolo.Debug("Init listening service for " + host + ":" + strconv.Itoa(port))

	ls.host = host
	ls.port = strconv.Itoa(port)
	ls.connectionChannel = connectionChannel

	ls.listenTo()

	ls.initialized = true
}

// Run listens to the port of this service and will start the handler.
func (ls *Listener) Run() {
	if !ls.initialized {
		sigolo.Fatal("Listening Service not initialized!")
	}

	for {
		conn, err := ls.waitForConnection()

		if err == nil {
			ls.connectionChannel(conn)
		} else {
			sigolo.Error(err.Error())
		}
	}
}

// listenTo actually listens to the port on the given host. It'll also exits the application if there's any problem.
func (ls *Listener) listenTo() {
	sigolo.Debug("Try to listen on port " + ls.port)

	listener, err := net.Listen("tcp", ls.host+":"+ls.port)

	if err == nil && listener != nil {
		sigolo.Info("Listen on " + ls.host + ":" + ls.port)
		ls.listener = listener
	} else if err != nil {
		sigolo.Error(err.Error())
		sigolo.Fatal("Maybe the port is not free?")
	} else if listener == nil {
		sigolo.Fatal("Could not listen to " + ls.host + ":" + ls.port + ". Unfortunately there's no error I could print here :( Check if no other services are running on port " + ls.port + ".")
	}
}

// waitForConnection accepts an incoming connection request.
func (ls *Listener) waitForConnection() (*net.Conn, error) {
	conn, err := ls.listener.Accept()

	if err == nil {
		sigolo.Info("Got connection on " + ls.host + ":" + ls.port)
		return &conn, nil
	}

	sigolo.Error(err.Error())
	return nil, err
}
