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

package local

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// The following statements are not required. They are type assertions to indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors at this location.

// They would normally be included in the commandController.go file, but they're located here for instructive purposes.
var _ = (infer.CustomResource[CommandInputs, CommandOutputs])((*Command)(nil))
var _ = (infer.CustomUpdate[CommandInputs, CommandOutputs])((*Command)(nil))
var _ = (infer.CustomDelete[CommandOutputs])((*Command)(nil))
var _ = (infer.ExplicitDependencies[CommandInputs, CommandOutputs])((*Command)(nil))

// This is the Create method. This will be run on every Command resource creation.
func (c *Command) Create(ctx p.Context, name string, input CommandInputs, preview bool) (string, CommandOutputs, error) {
	state := CommandOutputs{CommandInputs: input}
	id, err := resource.NewUniqueHex(name, 8, 0)
	if err != nil {
		return id, state, err
	}

	// If in preview, don't run the command.
	if preview {
		return id, state, nil
	}
	if input.Create == nil {
		return id, state, nil
	}
	cmd := *input.Create
	state.Stdout, state.Stderr, err = state.run(ctx, cmd)
	return id, state, err
}

// WireDependencies is relevant to secrets handling. This method indicates which Inputs
// the Outputs are derived from. If an output is derived from a secret input, the output
// will be a secret.

// This naive implementation conveys that every output is derived from all inputs.
func (r *Command) WireDependencies(f infer.FieldSelector, args *CommandInputs, state *CommandOutputs) {

	// get BaseInputs
	interpreterInput := f.InputField(&args.Interpreter)
	dirInput := f.InputField(&args.Dir)
	environmentInput := f.InputField(&args.Environment)
	stdinInput := f.InputField(&args.Stdin)
	assetPathsInput := f.InputField(&args.AssetPaths)
	archivePathsInput := f.InputField(&args.ArchivePaths)

	// set BaseInputs
	f.OutputField(&state.Interpreter).DependsOn(interpreterInput)
	f.OutputField(&state.Dir).DependsOn(dirInput)
	f.OutputField(&state.Environment).DependsOn(environmentInput)
	f.OutputField(&state.Stdin).DependsOn(stdinInput)
	f.OutputField(&state.AssetPaths).DependsOn(assetPathsInput)
	f.OutputField(&state.ArchivePaths).DependsOn(archivePathsInput)

	// get CommandInputs
	triggersInput := f.InputField(&args.Triggers)
	createInput := f.InputField(&args.Create)
	updateInput := f.InputField(&args.Update)
	deleteInput := f.InputField(&args.Delete)

	// set CommandInputs
	f.OutputField(&state.Triggers).DependsOn(triggersInput)
	f.OutputField(&state.Create).DependsOn(createInput)
	f.OutputField(&state.Update).DependsOn(updateInput)
	f.OutputField(&state.Delete).DependsOn(deleteInput)

	// Mark Stdout as derived from all inputs
	f.OutputField(&state.Stdout).DependsOn(
		triggersInput,
		createInput,
		updateInput,
		deleteInput,
		interpreterInput,
		dirInput,
		environmentInput,
		stdinInput,
		assetPathsInput,
		archivePathsInput,
	)

	// Mark Stderr as derived from all inputs
	f.OutputField(&state.Stderr).DependsOn(
		triggersInput,
		createInput,
		updateInput,
		deleteInput,
		interpreterInput,
		dirInput,
		environmentInput,
		stdinInput,
		assetPathsInput,
		archivePathsInput,
	)
}

// The Update method will be run on every update.
func (c *Command) Update(ctx p.Context, id string, olds CommandOutputs, news CommandInputs, preview bool) (CommandOutputs, error) {
	state := CommandOutputs{CommandInputs: news, BaseOutputs: olds.BaseOutputs}
	// If in preview, don't run the command.
	if preview {
		return state, nil
	}
	// Use Create command if Update is unspecified.
	cmd := news.Create
	if news.Update != nil {
		cmd = news.Update
	}
	// If neither are specified, do nothing.
	if cmd == nil {
		return state, nil
	}
	var err error
	state.Stdout, state.Stderr, err = state.run(ctx, *cmd)
	return state, err
}

// The Delete method will run when the resource is deleted.
func (c *Command) Delete(ctx p.Context, id string, props CommandOutputs) error {
	if props.Delete == nil {
		return nil
	}
	_, _, err := props.run(ctx, *props.Delete)
	return err
}
