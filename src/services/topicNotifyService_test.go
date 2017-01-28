package services_test

import (
	"net"
	"os"
	"testing"
)

var conn1, conn2 *net.Conn

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
