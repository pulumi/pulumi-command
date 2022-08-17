package remote

import (
	"fmt"
	"net"
	"time"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/retry"
	"golang.org/x/crypto/ssh"
)

type Connection struct {
	User       string   `pulumi:"user,optional"`
	Password   *string  `pulumi:"password,optional"`
	Host       string   `pulumi:"host"`
	Port       *float64 `pulumi:"port,optional"`
	PrivateKey *string  `pulumi:"privateKey,optional"`
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
}

func (con Connection) SShConfig() (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            con.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if con.PrivateKey != nil {
		signer, err := ssh.ParsePrivateKey([]byte(*con.PrivateKey))
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

	return config, nil
}

// Dial a ssh client connection from a ssh client configuration, retrying as necessary.
func (con Connection) Dial(ctx p.Context, config *ssh.ClientConfig) (*ssh.Client, error) {
	var client *ssh.Client
	var err error
	_, _, err = retry.Until(ctx, retry.Acceptor{
		Accept: func(try int, nextRetryTime time.Duration) (bool, interface{}, error) {
			client, err = ssh.Dial("tcp",
				net.JoinHostPort(con.Host, fmt.Sprintf("%d", con.Port)),
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
