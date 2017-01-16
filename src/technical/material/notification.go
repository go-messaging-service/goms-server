package material

import (
	"net"
)

type Notification struct {
	Connections *[]*net.Conn
	Data        string
}
