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
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.CustomResource[CommandInputs, CommandOutputs])((*Command)(nil))
var _ = (infer.CustomUpdate[CommandInputs, CommandOutputs])((*Command)(nil))
var _ = (infer.CustomDelete[CommandOutputs])((*Command)(nil))
var _ = (infer.ExplicitDependencies[CommandInputs, CommandOutputs])((*Command)(nil))

// WireDependencies is relevant to secrets handling. This method indicates which Inputs
// the Outputs are derived from. If an output is derived from a secret input, the output
// will be a secret.

// This naive implementation conveys that every output is derived from all inputs.
func (r *Command) WireDependencies(f infer.FieldSelector, args *CommandInputs, state *CommandOutputs) {
	createInput := f.InputField(&args.Create)
	updateInput := f.InputField(&args.Update)

	f.OutputField(&state.Connection).DependsOn(f.InputField(&args.Connection))
	f.OutputField(&state.Environment).DependsOn(f.InputField(&args.Environment))
	f.OutputField(&state.Triggers).DependsOn(f.InputField(&args.Triggers))
	f.OutputField(&state.Create).DependsOn(f.InputField(&args.Create))
	f.OutputField(&state.Delete).DependsOn(f.InputField(&args.Delete))
	f.OutputField(&state.Update).DependsOn(f.InputField(&args.Update))
	f.OutputField(&state.Stdin).DependsOn(f.InputField(&args.Stdin))

	f.OutputField(&state.Stdout).DependsOn(
		createInput,
		updateInput,
	)
	f.OutputField(&state.Stderr).DependsOn(
		createInput,
		updateInput,
	)
}

// This is the Create method. This will be run on every Command resource creation.
func (*Command) Create(ctx p.Context, name string, input CommandInputs, preview bool) (string, CommandOutputs, error) {
	state := CommandOutputs{CommandInputs: input}
	var err error
	id, err := resource.NewUniqueHex(name, 8, 0)
	if err != nil {
		return "", state, err
	}
	if preview {
		return id, state, nil
	}

	if state.Create == nil {
		return id, state, nil
	}
	cmd := ""
	if state.Create != nil {
		cmd = *state.Create
	}

	if !preview {
		err = state.run(ctx, cmd)
	}
	return id, state, err
}

// The Update method will be run on every update.
func (*Command) Update(ctx p.Context, id string, olds CommandOutputs, news CommandInputs, preview bool) (CommandOutputs, error) {
	state := CommandOutputs{CommandInputs: news, BaseOutputs: olds.BaseOutputs}
	if preview {
		return state, nil
	}
	var err error
	if !preview {
		if news.Update != nil {
			err = state.run(ctx, *news.Update)
		} else if news.Create != nil {
			err = state.run(ctx, *news.Create)
		}
	}
	return state, err
}

// The Delete method will run when the resource is deleted.
func (*Command) Delete(ctx p.Context, id string, props CommandOutputs) error {
	if props.Delete == nil {
		return nil
	}
	return props.run(ctx, *props.Delete)
}
