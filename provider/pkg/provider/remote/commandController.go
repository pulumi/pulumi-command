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
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.CustomResource[CommandInputs, CommandOutputs])((*Command)(nil))
var _ = (infer.CustomUpdate[CommandInputs, CommandOutputs])((*Command)(nil))
var _ = (infer.CustomDelete[CommandOutputs])((*Command)(nil))

// This is the Create method. This will be run on every Command resource creation.
func (*Command) Create(
	ctx context.Context,
	req infer.CreateRequest[CommandInputs],
) (infer.CreateResponse[CommandOutputs], error) {
	name := req.Name
	input := req.Inputs
	preview := req.DryRun
	state := CommandOutputs{CommandInputs: input}
	var err error
	id, err := resource.NewUniqueHex(name, 8, 0)
	if err != nil {
		return infer.CreateResponse[CommandOutputs]{ID: "", Output: state}, err
	}
	if preview {
		return infer.CreateResponse[CommandOutputs]{ID: id, Output: state}, nil
	}

	if input.Create == nil {
		return infer.CreateResponse[CommandOutputs]{ID: id, Output: state}, nil
	}
	cmd := ""
	if input.Create != nil {
		cmd = *input.Create
	}

	if !preview {
		err = state.run(ctx, cmd, input.Logging)
	}
	return infer.CreateResponse[CommandOutputs]{ID: id, Output: state}, err
}

// The Update method will be run on every update.
func (*Command) Update(
	ctx context.Context,
	req infer.UpdateRequest[CommandInputs, CommandOutputs],
) (infer.UpdateResponse[CommandOutputs], error) {
	olds := req.State
	news := req.Inputs
	preview := req.DryRun
	state := CommandOutputs{CommandInputs: news, BaseOutputs: olds.BaseOutputs}
	if preview {
		return infer.UpdateResponse[CommandOutputs]{Output: state}, nil
	}
	var err error
	if !preview {
		if news.Update != nil {
			err = state.run(ctx, *news.Update, news.Logging)
		} else if news.Create != nil {
			err = state.run(ctx, *news.Create, news.Logging)
		}
	}
	return infer.UpdateResponse[CommandOutputs]{Output: state}, err
}

// The Delete method will run when the resource is deleted.
func (*Command) Delete(ctx context.Context, req infer.DeleteRequest[CommandOutputs]) (infer.DeleteResponse, error) {
	props := req.State
	if props.Delete == nil {
		return infer.DeleteResponse{}, nil
	}
	return infer.DeleteResponse{}, props.run(ctx, *props.Delete, props.Logging)
}
