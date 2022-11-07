package local

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

type Run struct{}

func (r *Run) Annotate(a infer.Annotator) {
	a.Describe(&r, "A local command to be executed.\n"+
		"This command will always be run on any preview or deployment. "+
		"Use `local.Command` to avoid duplicating executions.")
}

type RunArgs struct {
	BaseArgs
	Command string `pulumi:"command"`
}

type RunState struct {
	RunArgs
	BaseState
}

func (*Run) Call(ctx p.Context, input RunArgs) (RunState, error) {
	r := RunState{RunArgs: input}
	var err error
	state := &CommandState{
		CommandArgs: CommandArgs{
			BaseArgs: input.BaseArgs,
		},
	}
	r.Stdout, r.Stderr, _, err = (state).run(ctx, input.Command)
	r.BaseState = state.BaseState
	return r, err
}
