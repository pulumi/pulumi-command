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
	"bytes"
	"context"
	"testing"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/common"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

// TestContext is a test implementation of p.Context that records all log messages in a buffer, regardless of severity.
type TestContext struct {
	context.Context
	output bytes.Buffer
}

func (c *TestContext) log(msg string) {
	c.output.WriteString(msg)
}

func (c *TestContext) Log(_ diag.Severity, msg string)                  { c.log(msg) }
func (c *TestContext) Logf(_ diag.Severity, msg string, _ ...any)       { c.log(msg) }
func (c *TestContext) LogStatus(_ diag.Severity, msg string)            { c.log(msg) }
func (c *TestContext) LogStatusf(_ diag.Severity, msg string, _ ...any) { c.log(msg) }
func (c *TestContext) RuntimeInformation() p.RunInfo                    { return p.RunInfo{} }

func TestOptionalLogging(t *testing.T) {
	for name, testCase := range map[string]struct {
		shouldLogOutput bool
		expectedLog     string
	}{
		"should log":     {shouldLogOutput: true, expectedLog: "hello"},
		"should not log": {shouldLogOutput: false, expectedLog: ""},
	} {
		t.Run(name, func(t *testing.T) {
			cmd := Command{}

			ctx := TestContext{Context: context.Background()}
			input := CommandInputs{
				BaseInputs: BaseInputs{
					CommonInputs: common.CommonInputs{
						LogOutput: pulumi.BoolRef(testCase.shouldLogOutput),
					},
				},
				ResourceInputs: common.ResourceInputs{
					Create: pulumi.StringRef("echo hello"),
				},
			}

			_, _, err := cmd.Create(&ctx, "name", input, false /* preview */)
			require.NoError(t, err)

			require.Equal(t, testCase.expectedLog, ctx.output.String())
		})
	}
}
