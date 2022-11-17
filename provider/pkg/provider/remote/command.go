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
type CommandArgs struct {
	// these field annotations allow you to
	// pulumi:"connection" specifies the name of the field in the schema
	// pulumi:"optional" specifies that a field is optional. This must be optional.
	// provider:"replaceOnChanges" specifies that a resource will be marked replaceOnChanges.
	// provider:"secret" specifies that a resource will be marked replaceOnChanges.
	Connection  *Connection        `pulumi:"connection,optional" provider:"secret"`
	Environment *map[string]string `pulumi:"environment,optional"`
	Triggers    *[]any             `pulumi:"triggers,optional"`
	Create      *string            `pulumi:"create,optional"`
	Delete      *string            `pulumi:"delete,optional"`
	Update      *string            `pulumi:"update,optional"`
	Stdin       *string            `pulumi:"stdin,optional"`
}

func (c *CommandArgs) Annotate(a infer.Annotator) {
	a.Describe(&c.Connection, "The parameters with which to connect to the remote host.")
	a.Describe(&c.Environment, "Additional environment variables available to the command's process.")
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.Create, "The command to run on create.")
	a.Describe(&c.Delete, "The command to run on delete.")
	a.Describe(&c.Update, "The command to run on update, if empty, create will run again.")
	a.Describe(&c.Stdin, "Pass a string to the command's process as standard in")
}

type BaseState struct {
	Stdout string `pulumi:"stdout"`
	Stderr string `pulumi:"stderr"`
}

func (c *BaseState) Annotate(a infer.Annotator) {
	a.Describe(&c.Stdout, "The standard output of the command's process")
	a.Describe(&c.Stderr, "The standard error of the command's process")
}

type CommandState struct {
	CommandArgs
	BaseState
}

func (c *CommandState) Annotate(a infer.Annotator) {
	a.Describe(&c.Stdout, "The standard output of the command's process")
	a.Describe(&c.Stderr, "The standard error of the command's process")
}
