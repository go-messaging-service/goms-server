package dist

import (
	"net"
)

type Notification struct {
	Connections *[]*net.Conn
	Topic       string
	Data        string
}
