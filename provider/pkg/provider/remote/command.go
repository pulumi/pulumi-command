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
)

type Command struct{}

func (c *Command) Annotate(a infer.Annotator) {
	a.Describe(&c, `A command to run on a remote host.
The connection is established via ssh.`)
}

// The arguments for a remote Command resource
type CommandInputs struct {
	// these field annotations allow you to
	// pulumi:"connection" specifies the name of the field in the schema
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	// provider:"secret" specifies that a field should be marked secret.
	Connection  *Connection        `pulumi:"connection" provider:"secret"`
	Environment *map[string]string `pulumi:"environment,optional"`
	Triggers    *[]any             `pulumi:"triggers,optional" provider:"replaceOnChanges"`
	Create      *string            `pulumi:"create,optional"`
	Delete      *string            `pulumi:"delete,optional"`
	Update      *string            `pulumi:"update,optional"`
	Stdin       *string            `pulumi:"stdin,optional"`
}

func (c *CommandInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Connection, "The parameters with which to connect to the remote host.")
	a.Describe(&c.Environment, "Additional environment variables available to the command's process.")
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.Create, "The command to run on create.")
	a.Describe(&c.Delete, "The command to run on delete.")
	a.Describe(&c.Update, "The command to run on update, if empty, create will run again.")
	a.Describe(&c.Stdin, "Pass a string to the command's process as standard in")
}

type BaseOutputs struct {
	Stdout string `pulumi:"stdout"`
	Stderr string `pulumi:"stderr"`
}

func (c *BaseOutputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Stdout, "The standard output of the command's process")
	a.Describe(&c.Stderr, "The standard error of the command's process")
}

type CommandOutputs struct {
	CommandInputs
	BaseOutputs
}

func (c *CommandOutputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Stdout, "The standard output of the command's process")
	a.Describe(&c.Stderr, "The standard error of the command's process")
}
