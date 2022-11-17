package local

import (
	"github.com/pulumi/pulumi-go-provider/infer"
)

type Run struct{}

func (r *RunArgs) Annotate(a infer.Annotator) {
	a.Describe(&r.Command, "The command to run.")
}

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
	// Marked as optional to match the old schema.
	// See https://github.com/pulumi/pulumi-command/issues/159
	Stdin  *string `pulumi:"stdin"`
	Stdout *string `pulumi:"stdout,optional"`
	BaseState
}
