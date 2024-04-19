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
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/common"
)

type Command struct{}

// Implementing Annotate lets you provide descriptions for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (c *Command) Annotate(a infer.Annotator) {
	a.Describe(&c, `A command to run on a remote host.
The connection is established via ssh.`)
}

// The arguments for a remote Command resource.
type CommandInputs struct {
	common.ResourceInputs
	common.CommonInputs
	// the pulumi-go-provider library uses field tags to dictate behavior.
	// pulumi:"connection" specifies the name of the field in the schema
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	// provider:"secret" specifies that a field should be marked secret.
	Connection  *Connection       `pulumi:"connection" provider:"secret"`
	Environment map[string]string `pulumi:"environment,optional"`
}

// Implementing Annotate lets you provide descriptions and default values for arguments and they will
// be visible in the provider's schema and the generated SDKs.
func (c *CommandInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Connection, "The parameters with which to connect to the remote host.")
	a.Describe(&c.Environment, `Additional environment variables available to the command's process.
Note that this only works if the SSH server is configured to accept these variables via AcceptEnv.
Alternatively, if a Bash-like shell runs the command on the remote host, you could prefix the command itself
with the variables in the form 'VAR=value command'.`)
}

// The properties for a remote Command resource.
type CommandOutputs struct {
	CommandInputs
	BaseOutputs
}
