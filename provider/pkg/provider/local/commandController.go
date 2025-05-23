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
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// The following statements are not required. They are type assertions to indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors at this location.

// They would normally be included in the commandController.go file, but they're located here for instructive purposes.
var _ = (infer.CustomResource[CommandInputs, CommandOutputs])((*Command)(nil))
var _ = (infer.CustomUpdate[CommandInputs, CommandOutputs])((*Command)(nil))
var _ = (infer.CustomDelete[CommandOutputs])((*Command)(nil))

// This is the Create method. This will be run on every Command resource creation.
func (c *Command) Create(ctx context.Context, req infer.CreateRequest[CommandInputs]) (infer.CreateResponse[CommandOutputs], error) {
	name := req.Name
	input := req.Inputs
	preview := req.DryRun
	state := CommandOutputs{CommandInputs: input}
	id, err := resource.NewUniqueHex(name, 8, 0)
	if err != nil {
		return infer.CreateResponse[CommandOutputs]{ID: id, Output: state}, err
	}

	// If in preview, don't run the command.
	if preview {
		return infer.CreateResponse[CommandOutputs]{ID: id, Output: state}, nil
	}
	if input.Create == nil {
		return infer.CreateResponse[CommandOutputs]{ID: id, Output: state}, nil
	}
	cmd := *input.Create
	err = run(ctx, cmd, state.BaseInputs, &state.BaseOutputs, input.Logging)
	return infer.CreateResponse[CommandOutputs]{ID: id, Output: state}, err
}

// WireDependencies controls how secrets and unknowns flow through a resource.
//
//	var _ = (infer.ExplicitDependencies[CommandInputs, CommandOutputs])((*Command)(nil))
//	func (r *Command) WireDependencies(f infer.FieldSelector, args *CommandInputs, state *CommandOutputs) { .. }
//
// Because we want every output to depend on every input, we can leave the default behavior.

// The Update method will be run on every update.
func (c *Command) Update(ctx context.Context, req infer.UpdateRequest[CommandInputs, CommandOutputs]) (infer.UpdateResponse[CommandOutputs], error) {
	olds := req.State
	news := req.Inputs
	preview := req.DryRun
	state := CommandOutputs{CommandInputs: news, BaseOutputs: olds.BaseOutputs}
	// If in preview, don't run the command.
	if preview {
		return infer.UpdateResponse[CommandOutputs]{Output: state}, nil
	}
	// Use Create command if Update is unspecified.
	cmd := news.Create
	if news.Update != nil {
		cmd = news.Update
	}
	// If neither are specified, do nothing.
	if cmd == nil {
		return infer.UpdateResponse[CommandOutputs]{Output: state}, nil
	}
	err := run(ctx, *cmd, state.BaseInputs, &state.BaseOutputs, news.Logging)
	return infer.UpdateResponse[CommandOutputs]{Output: state}, err
}

// The Delete method will run when the resource is deleted.
func (c *Command) Delete(ctx context.Context, req infer.DeleteRequest[CommandOutputs]) (infer.DeleteResponse, error) {
	props := req.State
	if props.Delete == nil {
		return infer.DeleteResponse{}, nil
	}
	return infer.DeleteResponse{}, run(ctx, *props.Delete, props.BaseInputs, &props.BaseOutputs, props.Logging)
}
