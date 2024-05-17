// Copyright 2024, Pulumi Corporation.
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
	"io"
	"strings"
	"testing"

	"github.com/gliderlabs/ssh"
	"github.com/pulumi/pulumi-command/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-command/provider/pkg/provider/local"
	"github.com/pulumi/pulumi-command/provider/pkg/provider/util/testutil"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestOptionalLogging(t *testing.T) {
	const (
		host = "localhost"
		port = 2222
	)

	// This SSH server always writes "foo" to stdout and "bar" to stderr, no matter the command.
	server := &ssh.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		Handler: func(s ssh.Session) {
			_, err := io.WriteString(s, "foo")
			require.NoError(t, err)
			_, err = io.WriteString(s.Stderr(), "bar")
			require.NoError(t, err)
		},
	}
	go func() {
		require.NoError(t, server.ListenAndServe())
	}()
	t.Cleanup(func() {
		_ = server.Close()
	})

	for _, logMode := range local.Logging.Values(local.LogStdoutAndStderr) {
		t.Run(logMode.Name, func(t *testing.T) {
			cmd := Command{}

			ctx := testutil.TestContext{Context: context.Background()}
			input := CommandInputs{
				CommonInputs: local.CommonInputs{
					Logging: &logMode.Value,
				},
				ResourceInputs: common.ResourceInputs{
					Create: pulumi.StringRef("ignored"),
				},
				Connection: &Connection{
					connectionBase: connectionBase{
						Host:           pulumi.StringRef(host),
						Port:           pulumi.Float64Ref(float64(port)),
						User:           pulumi.StringRef("user"), // unused but prevents nil panic
						PerDialTimeout: pulumi.IntRef(1),         // unused but prevents nil panic
					},
				},
			}

			_, _, err := cmd.Create(&ctx, "name", input, false /* preview */)
			require.NoError(t, err)

			log := ctx.Output.String()

			// When logging both stdout and stderr, the output could be foobar or barfoo.
			require.Equal(t, logMode.Value.ShouldLogStdout(), strings.Contains(log, "foo"))
			require.Equal(t, logMode.Value.ShouldLogStderr(), strings.Contains(log, "bar"))
		})
	}
}
