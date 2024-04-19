package common

import "github.com/pulumi/pulumi-go-provider/infer"

// ResourceInputs are inputs common to resource CRUD operations.
type ResourceInputs struct {
	// The field tags are used to provide metadata on the schema representation.
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	Triggers *[]any `pulumi:"triggers,optional" provider:"replaceOnChanges"`

	Create *string `pulumi:"create,optional"`
	Update *string `pulumi:"update,optional"`
	Delete *string `pulumi:"delete,optional"`
}

// CommonInputs are inputs common to each command execution, both resource and invoke.
type CommonInputs struct {
	// The field tags are used to provide metadata on the schema representation.
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	Stdin   *string  `pulumi:"stdin,optional"`
	Logging *Logging `pulumi:"logging,optional"`
}

type Logging string

const (
	LogStdout          Logging = "stdout"
	LogStderr          Logging = "stderr"
	LogStdoutAndStderr Logging = "stdoutAndStderr"
	NoLogging          Logging = "none"
)

func (Logging) Values() []infer.EnumValue[Logging] {
	return []infer.EnumValue[Logging]{
		{Name: string(LogStdout), Value: LogStdout, Description: "Capture stdout in logs but not stderr"},
		{Name: string(LogStderr), Value: LogStderr, Description: "Capture stderr in logs but not stdout"},
		{Name: string(LogStdoutAndStderr), Value: LogStdoutAndStderr, Description: "Capture stdout and stderr in logs"},
		{Name: string(NoLogging), Value: NoLogging, Description: "Capture no logs"},
	}
}

func (l *Logging) ShouldLogStdout() bool {
	return l == nil || *l == LogStdout || *l == LogStdoutAndStderr
}
func (l *Logging) ShouldLogStderr() bool {
	return l == nil || *l == LogStderr || *l == LogStdoutAndStderr
}

// Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *ResourceInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.Create, "The command to run on create.")
	a.Describe(&c.Delete, `The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
Command resource from previous create or update steps.`)
	a.Describe(&c.Update, `The command to run on update, if empty, create will 
run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
are set to the stdout and stderr properties of the Command resource from previous 
create or update steps.`)
}

func (c *CommonInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Stdin, "Pass a string to the command's process as standard in")
	a.Describe(&c.Logging, `If the command's stdout and stderr should be logged. This doesn't affect the capturing of
stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.`)
}
