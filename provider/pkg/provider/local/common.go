package local

// TODO This file (along with common_test.go) should be in the `common` package since its contents are used by `local`` and `remote`.
// It's in `local` for the time being due to pulumi/pulumi#16221.

import "github.com/pulumi/pulumi-go-provider/infer"

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

func (c *CommonInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Stdin, "Pass a string to the command's process as standard in")
	a.Describe(&c.Logging, `If the command's stdout and stderr should be logged. This doesn't affect the capturing of
stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.`)
}
