package test

import (
	"bufio"
	"net"
)

func initPipe() (*net.Conn, *bufio.Reader, *net.Conn, *bufio.Reader) {
	client, server := net.Pipe()

	clientReader := bufio.NewReader(client)
	serverReader := bufio.NewReader(server)

	return &client, clientReader, &server, serverReader
}
