package testutil

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/require"
)

type TestSshServer struct {
	Host string
	Port int64
}

// NewTestSshServer creates a new in-process SSH server with the specified handler.
// The server is bound to an arbitrary free port, and automatically closed
// during test cleanup.
func NewTestSshServer(t *testing.T, handler ssh.Handler) TestSshServer {
	const host = "127.0.0.1"

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, 0))
	require.NoErrorf(t, err, "net.Listen()")

	port, err := strconv.ParseInt(strings.Split(listener.Addr().String(), ":")[1], 10, 64)
	require.NoErrorf(t, err, "parse address %s allocated port number as int", listener.Addr())

	server := ssh.Server{Handler: handler}
	go func() {
		// "Serve always returns a non-nil error."
		_ = server.Serve(listener)
	}()
	t.Cleanup(func() {
		_ = server.Close()
	})

	return TestSshServer{
		Host: host,
		Port: port,
	}
}
