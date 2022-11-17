package remote

import (
	"bytes"
	"io"
	"os"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/util"
)

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.CustomResource[CommandArgs, CommandState])((*Command)(nil))
var _ = (infer.CustomUpdate[CommandArgs, CommandState])((*Command)(nil))
var _ = (infer.CustomDelete[CommandState])((*Command)(nil))
var _ = (infer.ExplicitDependencies[CommandArgs, CommandState])((*Command)(nil))

// WireDependencies marks the data dependencies between Inputs and Outputs
func (r *Command) WireDependencies(f infer.FieldSelector, args *CommandArgs, state *CommandState) {
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

func (*Command) Create(ctx p.Context, name string, input CommandArgs, preview bool) (string, CommandState, error) {
	state := CommandState{CommandArgs: input}
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
		state.Stdout, state.Stderr, err = state.run(ctx, cmd)
	}
	return id, state, err
}

func (*Command) Update(ctx p.Context, id string, olds CommandState, news CommandArgs, preview bool) (CommandState, error) {
	state := CommandState{CommandArgs: news}
	if preview {
		return state, nil
	}
	var err error
	if !preview {
		if news.Update != nil {
			state.Stdout, state.Stderr, err = state.run(ctx, *news.Update)
		} else if news.Create != nil {
			state.Stdout, state.Stderr, err = state.run(ctx, *news.Create)
		}
	}
	return state, err
}

func (*Command) Delete(ctx p.Context, id string, props CommandState) error {
	if props.Delete == nil {
		return nil
	}
	_, _, err := props.run(ctx, *props.Delete)
	return err
}

func (c *CommandState) run(ctx p.Context, cmd string) (string, string, error) {
	config, err := c.Connection.SShConfig()
	if err != nil {
		return "", "", err
	}

	client, err := c.Connection.Dial(ctx, config)
	if err != nil {
		return "", "", err
	}

	session, err := client.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()

	if c.Environment != nil {
		for k, v := range *c.Environment {
			session.Setenv(k, v)
		}
	}

	if c.Stdin != nil && len(*c.Stdin) > 0 {
		session.Stdin = strings.NewReader(*c.Stdin)
	}

	stdoutr, stdoutw, err := os.Pipe()
	if err != nil {
		return "", "", err
	}
	stderrr, stderrw, err := os.Pipe()
	if err != nil {
		return "", "", err
	}
	session.Stdout = stdoutw
	session.Stderr = stderrw

	var stdoutbuf bytes.Buffer
	var stderrbuf bytes.Buffer

	stdouttee := io.TeeReader(stdoutr, &stdoutbuf)
	stderrtee := io.TeeReader(stderrr, &stderrbuf)

	stdoutch := make(chan struct{})
	stderrch := make(chan struct{})
	go util.CopyOutput(ctx, stdouttee, stdoutch, diag.Debug)
	go util.CopyOutput(ctx, stderrtee, stderrch, diag.Info)

	err = session.Run(cmd)

	stdoutw.Close()
	stderrw.Close()

	<-stdoutch
	<-stderrch

	return stdoutbuf.String(), stderrbuf.String(), err
}
