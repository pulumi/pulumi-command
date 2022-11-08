package remote

import (
	"fmt"
	"net"
	"os"
	"time"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/retry"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	sshAgentSocketEnvVar = "SSH_AUTH_SOCK"
)

type Connection struct {
	User               string  `pulumi:"user,optional"`
	Password           *string `pulumi:"password,optional"`
	Host               string  `pulumi:"host"`
	Port               float64 `pulumi:"port,optional"`
	PrivateKey         *string `pulumi:"privateKey,optional"`
	PrivateKeyPassword *string `pulumi:"privateKeyPassword,optional"`
	AgentSocketPath    *string `pulumi:"agentSocketPath,optional"`
}

func (c *Connection) Annotate(a infer.Annotator) {
	a.Describe(&c, "Instructions for how to connect to a remote endpoint.")
	a.Describe(&c.User, "The user that we should use for the connection.")
	a.SetDefault(&c.User, "root")
	a.Describe(&c.Password, "The password we should use for the connection.")
	a.Describe(&c.Host, "The address of the resource to connect to.")
	a.Describe(&c.Port, "The port to connect to.")
	a.SetDefault(&c.Port, 22)
	a.Describe(&c.PrivateKey, "The contents of an SSH key to use for the connection. This takes preference over the password if provided.")
	a.Describe(&c.PrivateKeyPassword, "The password to use in case the private key is encrypted.")
	a.Describe(&c.AgentSocketPath, "SSH Agent socket path. Default to environment variable SSH_AUTH_SOCK if present.")
}

func (con *Connection) SShConfig() (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            con.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if con.PrivateKey != nil {
		var signer ssh.Signer
		var err error
		if con.PrivateKeyPassword != nil {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(*con.PrivateKey), []byte(*con.PrivateKeyPassword))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(*con.PrivateKey))
		}
		if err != nil {
			return nil, err
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}
	if con.Password != nil {
		config.Auth = append(config.Auth, ssh.Password(*con.Password))
		config.Auth = append(config.Auth, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
			for i := range questions {
				answers[i] = *con.Password
			}
			return answers, err
		}))
	}
	var sshAgentSocketPath *string
	if con.AgentSocketPath != nil {
		sshAgentSocketPath = con.AgentSocketPath
	}
	if envAgentSocketPath := os.Getenv(sshAgentSocketEnvVar); sshAgentSocketPath == nil && envAgentSocketPath != "" {
		sshAgentSocketPath = &envAgentSocketPath
	}
	if sshAgentSocketPath != nil {
		conn, err := net.Dial("unix", *sshAgentSocketPath)
		if err != nil {
			return nil, err
		}
		config.Auth = append(config.Auth, ssh.PublicKeysCallback(agent.NewClient(conn).Signers))
	}

	return config, nil
}

// Dial a ssh client connection from a ssh client configuration, retrying as necessary.
func (con *Connection) Dial(ctx p.Context, config *ssh.ClientConfig) (*ssh.Client, error) {
	var client *ssh.Client
	var err error
	_, _, err = retry.Until(ctx, retry.Acceptor{
		Accept: func(try int, nextRetryTime time.Duration) (bool, interface{}, error) {
			client, err = ssh.Dial("tcp",
				net.JoinHostPort(con.Host, fmt.Sprintf("%.0f", con.Port)),
				config)
			if err != nil {
				if try > 10 {
					return true, nil, err
				}
				return false, nil, nil
			}
			return true, nil, nil
		},
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
