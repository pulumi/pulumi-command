//go:build !windows

package remote

import (
	"net"
	"os"
)

const sshAgentSocketEnvVar = "SSH_AUTH_SOCK"

func tryGetDefaultAgentSocket() *string {
	if envAgentSocketPath := os.Getenv(sshAgentSocketEnvVar); envAgentSocketPath != "" {
		return &envAgentSocketPath
	}
	return nil
}

func dialAgent(path string) (net.Conn, error) {
	return net.Dial("unix", path)
}
