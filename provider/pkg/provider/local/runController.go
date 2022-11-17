package local

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// These are not required. They indicate to Go that Run implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.ExplicitDependencies[RunArgs, RunState])((*Run)(nil))

// WireDependencies marks the data dependencies between Inputs and Outputs
func (r *Run) WireDependencies(f infer.FieldSelector, args *RunArgs, state *RunState) {

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

func (*Run) Call(ctx p.Context, input RunArgs) (RunState, error) {
	r := RunState{RunArgs: input}
	var err error
	state := &CommandState{
		CommandArgs: CommandArgs{
			BaseArgs: input.BaseArgs,
		},
	}
	*r.Stdout, r.Stderr, err = (state).run(ctx, input.Command)
	r.BaseState = state.BaseState
	return r, err
}
