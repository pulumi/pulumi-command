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
	"bytes"
	"io"
	"os"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/util"
)

func (c *CommandOutputs) run(ctx p.Context, cmd string) (string, string, error) {
	config, err := c.Connection.SShConfig()
	if err != nil {
		return "", "", err
	}

	client, err := c.Connection.Dial(ctx, config)
	if err != nil {
		return "", "", err
	}

	session, err := client.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()

	if c.Environment != nil {
		for k, v := range *c.Environment {
			session.Setenv(k, v)
		}
	}
	if c.Stdout != "" {
		session.Setenv(util.PULUMI_COMMAND_STDOUT, c.Stdout)
	}
	if c.Stderr != "" {
		session.Setenv(util.PULUMI_COMMAND_STDERR, c.Stderr)
	}

	if c.Stdin != nil && len(*c.Stdin) > 0 {
		session.Stdin = strings.NewReader(*c.Stdin)
	}

	stdoutr, stdoutw, err := os.Pipe()
	if err != nil {
		return "", "", err
	}
	stderrr, stderrw, err := os.Pipe()
	if err != nil {
		return "", "", err
	}
	session.Stdout = stdoutw
	session.Stderr = stderrw

	var stdoutbuf bytes.Buffer
	var stderrbuf bytes.Buffer

	stdouttee := io.TeeReader(stdoutr, &stdoutbuf)
	stderrtee := io.TeeReader(stderrr, &stderrbuf)

	stdoutch := make(chan struct{})
	stderrch := make(chan struct{})
	go util.CopyOutput(ctx, stdouttee, stdoutch, diag.Debug)
	go util.CopyOutput(ctx, stderrtee, stderrch, diag.Info)

	err = session.Run(cmd)

	stdoutw.Close()
	stderrw.Close()

	<-stdoutch
	<-stderrch

	return stdoutbuf.String(), stderrbuf.String(), err
}
