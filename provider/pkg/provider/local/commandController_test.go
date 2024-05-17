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

package local

import (
	"context"
	"strings"
	"testing"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-command/provider/pkg/provider/util/testutil"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestOptionalLogging(t *testing.T) {
	for _, logMode := range Logging.Values(LogStdoutAndStderr) {

		t.Run(logMode.Name, func(t *testing.T) {
			cmd := Command{}

			ctx := testutil.TestContext{Context: context.Background()}
			input := CommandInputs{
				BaseInputs: BaseInputs{
					Logging: &logMode.Value,
				},
				ResourceInputs: common.ResourceInputs{
					Create: pulumi.StringRef("echo foo; echo bar >> /dev/stderr"),
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
