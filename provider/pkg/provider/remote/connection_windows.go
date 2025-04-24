//go:build windows

package remote

import (
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

const windowsDefaultAgentPipe = `\\.\pipe\openssh-ssh-agent`

func tryGetDefaultAgentSocket() *string {
	defaultPipe := windowsDefaultAgentPipe
	timeout := 100 * time.Millisecond
	// Try to connect with a short timeout to check if the pipe exists
	if conn, err := winio.DialPipe(defaultPipe, &timeout); err == nil {
		conn.Close()
		return &defaultPipe
	}
	return nil
}

func dialAgent(path string) (net.Conn, error) {
	return winio.DialPipe(path, nil)
}
