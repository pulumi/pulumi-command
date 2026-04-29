// Package sshfixture starts ephemeral OpenSSH server containers for the
// integration tests. It exposes one host and an optional jump-host pair so
// the examples can exercise remote.Command and CopyToRemote without any
// cloud credentials.
package sshfixture

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/ssh"
)

// Pin to a specific upstream tag so CI is reproducible.
const opensshImage = "lscr.io/linuxserver/openssh-server:version-10.2_p1-r0"

// User is the SSH login on every fixture container.
const User = "testuser"

// Server describes one running OpenSSH container as seen from the host
// running the test.
type Server struct {
	Host          string
	Port          int
	PrivateKeyPEM string
}

// Proxy describes a jump-host scenario: the test connects to Proxy from the
// host, and the proxy reaches Target via a Docker network alias. Both servers
// share the same SSH key.
type Proxy struct {
	Proxy         Server
	Target        Server // Host is a Docker-network DNS name, Port is 2222.
	PrivateKeyPEM string
}

// keyPair generates a fresh ed25519 keypair and returns
// (private key in OpenSSH PEM, public key in authorized_keys format).
func keyPair(t *testing.T) (privPEM, pubAuthorized string) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pemBlock, err := ssh.MarshalPrivateKey(priv, "")
	require.NoError(t, err)
	privPEM = string(pem.EncodeToMemory(pemBlock))

	sshPub, err := ssh.NewPublicKey(pub)
	require.NoError(t, err)
	pubAuthorized = string(ssh.MarshalAuthorizedKey(sshPub))
	return privPEM, pubAuthorized
}

// startContainer launches an openssh-server container with the given public
// key authorized. Extra options (e.g. attaching to a network) are applied on
// top of the base request.
func startContainer(
	ctx context.Context, t *testing.T, pubKey string, extra ...testcontainers.CustomizeRequestOption,
) testcontainers.Container {
	t.Helper()

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        opensshImage,
			ExposedPorts: []string{"2222/tcp"},
			Env: map[string]string{
				"PUID":            "1000",
				"PGID":            "1000",
				"TZ":              "Etc/UTC",
				"USER_NAME":       User,
				"PUBLIC_KEY":      pubKey,
				"PASSWORD_ACCESS": "false",
				"SUDO_ACCESS":     "false",
			},
			Files: []testcontainers.ContainerFile{{
				// The image ships with AllowTcpForwarding disabled, which
				// breaks the proxy connection path. Drop an override into
				// sshd_config.d so the jump-host scenario works.
				Reader:            strings.NewReader("AllowTcpForwarding yes\n"),
				ContainerFilePath: "/config/sshd/sshd_config.d/allow-tcp-forwarding.conf",
				FileMode:          0o644,
			}},
			// The host port becomes connectable as soon as Docker proxies
			// it, which can be earlier than sshd actually accepting
			// connections. The linuxserver image prints this once init
			// finishes and sshd is fully up.
			WaitingFor: wait.ForAll(
				wait.ForLog("[ls.io-init] done.").
					WithStartupTimeout(90*time.Second),
				wait.ForListeningPort("2222/tcp").
					WithStartupTimeout(90*time.Second),
			),
		},
		Started: true,
	}
	for _, opt := range extra {
		require.NoError(t, opt.Customize(&req))
	}

	c, err := testcontainers.GenericContainer(ctx, req)
	testcontainers.CleanupContainer(t, c)
	require.NoError(t, err, "starting openssh-server container")
	return c
}

// New starts a single openssh-server container reachable from the host.
func New(t *testing.T) Server {
	t.Helper()
	ctx := context.Background()

	priv, pub := keyPair(t)
	c := startContainer(ctx, t, pub)

	host, err := c.Host(ctx)
	require.NoError(t, err)
	port, err := c.MappedPort(ctx, "2222/tcp")
	require.NoError(t, err)

	return Server{
		Host:          host,
		Port:          port.Int(),
		PrivateKeyPEM: priv,
	}
}

// NewProxy starts two openssh-server containers wired through a private
// Docker network so the second is only reachable via the first.
func NewProxy(t *testing.T) Proxy {
	t.Helper()
	ctx := context.Background()

	net, err := network.New(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { _ = net.Remove(ctx) })

	priv, pub := keyPair(t)

	const (
		proxyAlias  = "ssh-proxy"
		targetAlias = "ssh-target"
	)

	// target is reachable only via the network alias; we don't query its
	// mapped host port. Keep the reference alive so the cleanup hook the
	// testcontainers runtime registers stays bound to it.
	_ = startContainer(ctx, t, pub,
		network.WithNetwork([]string{targetAlias}, net),
	)
	proxy := startContainer(ctx, t, pub,
		network.WithNetwork([]string{proxyAlias}, net),
	)

	proxyHost, err := proxy.Host(ctx)
	require.NoError(t, err)
	proxyPort, err := proxy.MappedPort(ctx, "2222/tcp")
	require.NoError(t, err)

	return Proxy{
		Proxy: Server{
			Host:          proxyHost,
			Port:          proxyPort.Int(),
			PrivateKeyPEM: priv,
		},
		Target: Server{
			Host:          targetAlias,
			Port:          2222,
			PrivateKeyPEM: priv,
		},
		PrivateKeyPEM: priv,
	}
}

// Endpoint formats a host:port string for logging.
func (s Server) Endpoint() string { return fmt.Sprintf("%s:%d", s.Host, s.Port) }
