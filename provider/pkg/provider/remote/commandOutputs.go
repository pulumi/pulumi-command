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
	"fmt"
	"io"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/util"
)

func (c *CommandOutputs) run(ctx p.Context, cmd string) error {
	client, err := c.Connection.Dial(ctx)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
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

	var stdoutbuf, stderrbuf, stdouterrbuf bytes.Buffer
	r, w := io.Pipe()
	session.Stdout = io.MultiWriter(&stdoutbuf, &stdouterrbuf, w)
	session.Stderr = io.MultiWriter(&stderrbuf, &stdouterrbuf, w)

	stdouterrch := make(chan struct{})
	go util.CopyOutput(ctx, r, stdouterrch, diag.Info)

	err = session.Run(cmd)

	w.Close()
	<-stdouterrch

	if err != nil {
		return fmt.Errorf("%w: running %q:\n%s", err, cmd, stdouterrbuf.String())
	}
	c.BaseOutputs = BaseOutputs{
		Stdout: strings.TrimSuffix(stdoutbuf.String(), "\n"),
		Stderr: strings.TrimSuffix(stderrbuf.String(), "\n"),
	}
	return nil
}
