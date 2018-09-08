package conn

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

func (l *Listener) Init(host string, port int, topics []string, connectionChannel func(*net.Conn)) {
	sigolo.Debug("Init listening service for " + host + ":" + strconv.Itoa(port))

	l.host = host
	l.port = strconv.Itoa(port)
	l.connectionChannel = connectionChannel

	l.listenTo()

	l.initialized = true
}

// Run listens to the port of this service and will start the handler.
func (l *Listener) Run() {
	if !l.initialized {
		sigolo.Fatal("Listening Service not initialized!")
	}

	for {
		conn, err := l.waitForConnection()

		if err == nil {
			l.connectionChannel(conn)
		} else {
			sigolo.Error(err.Error())
		}
	}
}

// listenTo actually listens to the port on the given host. It'll also exits the application if there's any problem.
func (l *Listener) listenTo() {
	sigolo.Debug("Try to listen on port " + l.port)

	listener, err := net.Listen("tcp", l.host+":"+l.port)

	if err == nil && listener != nil {
		sigolo.Info("Listen on " + l.host + ":" + l.port)
		l.listener = listener
	} else if err != nil {
		sigolo.Error(err.Error())
		sigolo.Fatal("Maybe the port is not free?")
	} else if listener == nil {
		sigolo.Fatal("Could not listen to " + l.host + ":" + l.port + ". Unfortunately there's no error I could print here :( Check if no other services are running on port " + l.port + ".")
	}
}

// waitForConnection accepts an incoming connection request.
func (l *Listener) waitForConnection() (*net.Conn, error) {
	conn, err := l.listener.Accept()

	if err == nil {
		sigolo.Info("Got connection on " + l.host + ":" + l.port)
		return &conn, nil
	}

	sigolo.Error(err.Error())
	return nil, err
}
