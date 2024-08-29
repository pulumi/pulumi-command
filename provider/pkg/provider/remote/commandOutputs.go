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
	"context"
	"fmt"
	"io"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/util"
)

func (c *CommandOutputs) run(ctx context.Context, cmd string, logging *Logging) error {
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
		for k, v := range c.Environment {
			err := session.Setenv(k, v)
			if err != nil {
				return logAndWrapSetenvErr(diag.Error, k, ctx, err)
			}
		}
	}

	if c.AddPreviousOutputInEnv == nil || *c.AddPreviousOutputInEnv {
		// Set remote Stdout and Stderr environment variables optimistically, but log and continue if they fail.
		if c.Stdout != "" {
			err := session.Setenv(util.PULUMI_COMMAND_STDOUT, c.Stdout)
			if err != nil {
				// Set remote Stdout var optimistically, but warn and continue on failure.
				//
				//nolint:errcheck
				logAndWrapSetenvErr(diag.Warning, util.PULUMI_COMMAND_STDOUT, ctx, err)
			}
		}
		if c.Stderr != "" {
			err := session.Setenv(util.PULUMI_COMMAND_STDERR, c.Stderr)
			if err != nil {
				// Set remote STDERR var optimistically, but warn and continue on failure.
				//
				//nolint:errcheck
				logAndWrapSetenvErr(diag.Warning, util.PULUMI_COMMAND_STDERR, ctx, err)
			}
		}
	}

	if c.Stdin != nil && len(*c.Stdin) > 0 {
		session.Stdin = strings.NewReader(*c.Stdin)
	}

	var stdoutbuf, stderrbuf, stdouterrbuf bytes.Buffer
	r, w := io.Pipe()

	stdoutWriters := []io.Writer{&stdoutbuf, &stdouterrbuf}
	if logging.ShouldLogStdout() {
		stdoutWriters = append(stdoutWriters, w)
	}
	session.Stdout = io.MultiWriter(stdoutWriters...)

	stderrWriters := []io.Writer{&stderrbuf, &stdouterrbuf}
	if logging.ShouldLogStderr() {
		stderrWriters = append(stderrWriters, w)
	}
	session.Stderr = io.MultiWriter(stderrWriters...)

	stdouterrch := make(chan struct{})
	go util.LogOutput(ctx, r, stdouterrch, diag.Info)

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

func logAndWrapSetenvErr(severity diag.Severity, key string, ctx context.Context, err error) error {
	l := p.GetLogger(ctx)
	msg := fmt.Sprintf(`Unable to set '%s'. This only works if your SSH server is configured to accept
	these variables via AcceptEnv. Alternatively, if a Bash-like shell runs the command on the remote host, you could
	prefix the command itself with the variables in the form 'VAR=value command'`, key)
	switch severity {
	case diag.Error:
		l.Error(msg)
	case diag.Warning:
		l.Warning(msg)
	case diag.Info:
		l.Info(msg)
	default:
		l.Debug(msg)
	}
	return fmt.Errorf("could not set environment variable %q: %w", key, err)
}
