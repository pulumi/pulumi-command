// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package remote

import (
	"context"
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

var (
	dialErrorDefault   = 10
	dialErrorUnlimited = -1
)

type Connection struct {
	connectionBase
	Proxy *ProxyConnection `pulumi:"proxy,optional"`
}

type connectionBase struct {
	User               *string  `pulumi:"user,optional"`
	Password           *string  `pulumi:"password,optional"`
	Host               *string  `pulumi:"host"`
	Port               *float64 `pulumi:"port,optional"`
	PrivateKey         *string  `pulumi:"privateKey,optional"`
	PrivateKeyPassword *string  `pulumi:"privateKeyPassword,optional"`
	AgentSocketPath    *string  `pulumi:"agentSocketPath,optional"`
	DialErrorLimit     *int     `pulumi:"dialErrorLimit,optional"`
	PerDialTimeout     *int     `pulumi:"perDialTimeout,optional"`
}

func (c *Connection) Annotate(a infer.Annotator) {
	a.Describe(&c, "Instructions for how to connect to a remote endpoint.")
	a.Describe(&c.User, "The user that we should use for the connection.")
	a.SetDefault(&c.User, "root")
	a.Describe(&c.Password, "The password we should use for the connection.")
	a.Describe(&c.Host, "The address of the resource to connect to.")
	a.Describe(&c.Port, "The port to connect to. Defaults to 22.")
	a.SetDefault(&c.Port, 22)
	a.Describe(&c.PrivateKey, "The contents of an SSH key to use for the connection. This takes preference over the password if provided.")
	a.Describe(&c.PrivateKeyPassword, "The password to use in case the private key is encrypted.")
	a.Describe(&c.AgentSocketPath, "SSH Agent socket path. Default to environment variable SSH_AUTH_SOCK if present.")
	a.Describe(&c.DialErrorLimit, "Max allowed errors on trying to dial the remote host. -1 set count to unlimited. Default value is 10.")
	a.Describe(&c.Proxy, "The connection settings for the bastion/proxy host.")
	a.SetDefault(&c.DialErrorLimit, dialErrorDefault)
	a.Describe(&c.PerDialTimeout, "Max number of seconds for each dial attempt. 0 implies no maximum. Default value is 15 seconds.")
	a.SetDefault(&c.PerDialTimeout, 15)
}

func (con *connectionBase) SShConfig() (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            *con.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * time.Duration(*con.PerDialTimeout),
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
		config.Auth = append(config.Auth, ssh.KeyboardInteractive(
			func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				answers := make([]string, len(questions))
				for i := range questions {
					answers[i] = *con.Password
				}
				return answers, nil
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

func dialWithRetry[T any](ctx context.Context, msg string, maxAttempts int, f func() (T, error)) (T, error) {
	var userError error
	ok, data, err := retry.Until(ctx, retry.Acceptor{
		Accept: func(try int, _ time.Duration) (bool, any, error) {
			var result T
			result, userError = f()
			if userError == nil {
				return true, result, nil
			}
			dials := try + 1
			if reachedDialingErrorLimit(dials, maxAttempts) {
				return true, nil, fmt.Errorf("after %d failed attempts: %w",
					try, userError)
			}
			var limit string
			if maxAttempts == -1 {
				limit = "inf"
			} else {
				limit = fmt.Sprintf("%d", maxAttempts)
			}
			p.GetLogger(ctx).InfoStatusf("%s %d/%s failed: retrying",
				msg, dials, limit)
			return false, nil, nil
		},
	})
	// It's important to check both `ok` and `err` as sometimes `err` will be nil when `ok` is false, such as when the context is cancelled.
	if ok && err == nil {
		return data.(T), nil
	}

	var t T
	if err == nil {
		// `err` is nil but ok was false, use the err reported from the context.
		err = ctx.Err()
	}
	return t, err
}

// Dial a ssh client connection from a ssh client configuration, retrying as necessary.
func (con *Connection) Dial(ctx context.Context) (*ssh.Client, error) {
	config, err := con.SShConfig()
	if err != nil {
		return nil, err
	}

	endpoint := net.JoinHostPort(*con.Host, fmt.Sprintf("%d", int(*con.Port)))
	tries := con.getDialErrorLimit()
	if con.Proxy == nil {
		return dialWithRetry(ctx, "Dial", tries, func() (*ssh.Client, error) {
			return ssh.Dial("tcp", endpoint, config)
		})
	}

	proxyConfig, err := con.Proxy.SShConfig()
	if err != nil {
		return nil, fmt.Errorf("proxy: %w", err)
	}

	proxyTries := con.Proxy.getDialErrorLimit()
	// The user has specified a proxy connection. First, connect to the proxy:
	proxyClient, err := dialWithRetry(ctx, "Dial proxy", proxyTries, func() (*ssh.Client, error) {
		return ssh.Dial("tcp",
			net.JoinHostPort(*con.Proxy.Host, fmt.Sprintf("%d", int(*con.Proxy.Port))),
			proxyConfig)
	})
	if err != nil {
		return nil, fmt.Errorf("proxy: %w", err)
	}

	// Having connected with the proxy, we establish a connection from our proxy to
	// our server.
	conn, err := dialWithRetry(ctx, "Dial from proxy", tries, func() (net.Conn, error) {
		return proxyClient.Dial("tcp", endpoint)
	})
	if err != nil {
		return nil, fmt.Errorf("proxy: %w", err)
	}

	// We initiate a SSH connection over the bridge we just established.
	var channel <-chan ssh.NewChannel
	var req <-chan *ssh.Request
	proxyConn, err := dialWithRetry(ctx, "Dial", tries, func() (ssh.Conn, error) {
		c, ch, r, err := ssh.NewClientConn(conn, endpoint, config)
		channel = ch
		req = r
		return c, err
	})
	if err != nil {
		return nil, fmt.Errorf("proxy: %w", err)
	}

	return ssh.NewClient(proxyConn, channel, req), nil
}

func (con connectionBase) getDialErrorLimit() int {
	if con.DialErrorLimit == nil {
		return dialErrorDefault
	}
	return *con.DialErrorLimit
}

func reachedDialingErrorLimit(dials, dialErrorLimit int) bool {
	return dialErrorLimit > dialErrorUnlimited &&
		dials > dialErrorLimit
}
