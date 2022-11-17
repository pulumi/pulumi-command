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
)

// These are not required. They indicate to Go that Run implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.ExplicitDependencies[RunInputs, RunOutputs])((*Run)(nil))

// WireDependencies marks the data dependencies between Inputs and Outputs
func (r *Run) WireDependencies(f infer.FieldSelector, args *RunInputs, state *RunOutputs) {

	interpreterInput := f.InputField(&args.Interpreter)
	dirInput := f.InputField(&args.Dir)
	environmentInput := f.InputField(&args.Environment)
	stdinInput := f.InputField(&args.Stdin)
	assetPathsInput := f.InputField(&args.AssetPaths)
	archivePathsInput := f.InputField(&args.ArchivePaths)

	f.OutputField(&state.Interpreter).DependsOn(interpreterInput)
	f.OutputField(&state.Dir).DependsOn(dirInput)
	f.OutputField(&state.Environment).DependsOn(environmentInput)
	f.OutputField(&state.Stdin).DependsOn(stdinInput)
	f.OutputField(&state.AssetPaths).DependsOn(assetPathsInput)
	f.OutputField(&state.ArchivePaths).DependsOn(archivePathsInput)

	commandInput := f.InputField(&args.Command)

	f.OutputField(&state.Stdout).DependsOn(
		commandInput,
		interpreterInput,
		dirInput,
		environmentInput,
		stdinInput,
		assetPathsInput,
		archivePathsInput,
	)

	f.OutputField(&state.Stderr).DependsOn(
		commandInput,
		interpreterInput,
		dirInput,
		environmentInput,
		stdinInput,
		assetPathsInput,
		archivePathsInput,
	)
}

func (*Run) Call(ctx p.Context, input RunInputs) (RunOutputs, error) {
	r := RunOutputs{RunInputs: input}
	var err error
	state := &CommandOutputs{
		CommandInputs: CommandInputs{
			BaseInputs: input.BaseInputs,
		},
	}
	r.Stdout, r.Stderr, err = (state).run(ctx, input.Command)
	r.BaseOutputs = state.BaseOutputs
	return r, err
}
