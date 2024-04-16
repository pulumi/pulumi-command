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
	"testing"

	"github.com/gliderlabs/ssh"
	"github.com/pulumi/pulumi-command/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-command/provider/pkg/provider/util"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestOptionalLogging(t *testing.T) {
	const (
		host           = "localhost"
		port           = 2222
		serverResponse = "look, it's SSH!"
	)

	server := &ssh.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		Handler: func(s ssh.Session) {
			_, err := io.WriteString(s, serverResponse)
			require.NoError(t, err)
		},
	}
	go func() {
		require.NoError(t, server.ListenAndServe())
	}()
	t.Cleanup(func() {
		_ = server.Close()
	})

	for name, testCase := range map[string]struct {
		shouldLogOutput bool
		expectedLog     string
	}{
		"should log":     {shouldLogOutput: true, expectedLog: serverResponse},
		"should not log": {shouldLogOutput: false, expectedLog: ""},
	} {
		t.Run(name, func(t *testing.T) {
			cmd := Command{}

			ctx := util.TestContext{Context: context.Background()}
			input := CommandInputs{
				CommonInputs: common.CommonInputs{
					LogOutput: pulumi.BoolRef(testCase.shouldLogOutput),
				},
				ResourceInputs: common.ResourceInputs{
					Create: pulumi.StringRef("ignored"),
				},
				Connection: &Connection{
					connectionBase: connectionBase{
						Host:           pulumi.StringRef(host),
						Port:           pulumi.Float64Ref(float64(port)),
						User:           pulumi.StringRef("foo"), // unused but prevents nil panic
						PerDialTimeout: pulumi.IntRef(1),        // unused but prevents nil panic
					},
				},
			}

			_, _, err := cmd.Create(&ctx, "name", input, false /* preview */)
			require.NoError(t, err)

			require.Equal(t, testCase.expectedLog, ctx.Output.String())
		})
	}
}
