package main

import (
	"./logger"
	"net"
)

var listener net.Listener

func listenToPort(port string) {
	logger.Info("Try to listen on port " + port)

	ln, err := net.Listen("tcp", ":"+port)
	listener = ln

	if err == nil {
		logger.Info("Got listener for port " + port)
	} else {
		logger.Err(err.Error())
	}
}

func waitForConnection() (*net.Conn, error) {
	conn, err := listener.Accept()

	if err == nil {
		logger.Info("Got connection :D")
		return &conn, nil
	}

	logger.Err(err.Error())
	return nil, err
}
