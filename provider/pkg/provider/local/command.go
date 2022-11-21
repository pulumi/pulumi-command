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

// This file contains metadata around the types for

package local

import (
	"github.com/pulumi/pulumi-go-provider/infer"
)

// This is the type that implements the Command resource methods.
// The methods are declared in the commandController.go file.
type Command struct{}

// The following statement is not required. It is a type assertion to indicate to Go that Command
// implements the following interfaces. If the function signature doesn't match or isn't implemented,
// we get nice compile time errors at this location.

var _ = (infer.Annotated)((*Command)(nil))

// Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (c *Command) Annotate(a infer.Annotator) {
	a.Describe(&c, "A local command to be executed.\n"+
		"This command can be inserted into the life cycles of other resources using the\n"+
		"`dependsOn` or `parent` resource options. A command is considered to have\n"+
		"failed when it finished with a non-zero exit code. This will fail the CRUD step\n"+
		"of the `Command` resource.")
}

// These are the inputs (or arguments) to a Command resource.
type CommandInputs struct {
	// The field tags are used to provide metadata on the schema representation.
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	Triggers *[]any  `pulumi:"triggers,optional" provider:"replaceOnChanges"`
	Create   *string `pulumi:"create,optional"`
	Update   *string `pulumi:"update,optional"`
	Delete   *string `pulumi:"delete,optional"`
	BaseInputs
}

// Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *CommandInputs) Annotate(a infer.Annotator) {
	c.BaseInputs.Annotate(a)
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.Create, "The command to run on create.")
	a.Describe(&c.Delete, "The command to run on delete.")
	a.Describe(&c.Update, "The command to run on update, if empty, create will run again.")
}

// These are the outputs (or properties) of a Command resource.
type CommandOutputs struct {
	CommandInputs
	BaseOutputs
}
