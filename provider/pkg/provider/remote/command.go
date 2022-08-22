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

type Command struct{}

var _ = (infer.CustomDelete[CommandState])((*Command)(nil))
var _ = (infer.CustomUpdate[CommandArgs, CommandState])((*Command)(nil))

func (c *Command) Annotate(a infer.Annotator) {
	a.Describe(&c, `A command to run on a remote host.
The connection is established via ssh.`)
}

type CommandArgs struct {
	Connection  Connection         `pulumi:"connection" provider:"replaceOnChanges"`
	Environment *map[string]string `pulumi:"environment,optional"`
	Triggers    *[]any             `pulumi:"triggers,optional" provider:"replaceOnChanges"`
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

type CommandState struct {
	CommandArgs
	Stdout string `pulumi:"stdout"`
	Stderr string `pulumi:"stderr"`
}

func (c *CommandState) Annotate(a infer.Annotator) {
	a.Describe(&c.Stdout, "The standard output of the command's process")
	a.Describe(&c.Stderr, "The standard error of the command's process")
}

func (*Command) Create(ctx p.Context, name string, input CommandArgs, preview bool) (string, CommandState, error) {
	s := CommandState{CommandArgs: input}
	var id string
	var err error
	cmd := ""
	if s.Create != nil {
		cmd = *s.Create
	}
	if !preview {
		s.Stdout, s.Stderr, id, err = s.run(ctx, cmd)
	}
	return id, s, err
}

func (*Command) Update(ctx p.Context, id string, olds CommandState, news CommandArgs, preview bool) (CommandState, error) {
	state := CommandState{CommandArgs: news}
	if preview {
		return state, nil
	}
	var err error
	if !preview {
		if news.Update != nil {
			state.Stdout, state.Stderr, _, err = state.run(ctx, *news.Update)
		} else if news.Create != nil {
			state.Stdout, state.Stderr, _, err = state.run(ctx, *news.Create)
		}
	}
	return state, err
}

func (*Command) Delete(ctx p.Context, id string, props CommandState) error {
	if props.Delete == nil {
		return nil
	}
	_, _, _, err := props.run(ctx, *props.Delete)
	return err
}

func (c *CommandState) run(ctx p.Context, cmd string) (string, string, string, error) {
	config, err := c.Connection.SShConfig()
	if err != nil {
		return "", "", "", err
	}

	client, err := c.Connection.Dial(ctx, config)
	if err != nil {
		return "", "", "", err
	}

	session, err := client.NewSession()
	if err != nil {
		return "", "", "", err
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

	id, err := resource.NewUniqueHex("", 8, 0)
	if err != nil {
		return "", "", "", err
	}

	stdoutr, stdoutw, err := os.Pipe()
	if err != nil {
		return "", "", "", err
	}
	stderrr, stderrw, err := os.Pipe()
	if err != nil {
		return "", "", "", err
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

	return stdoutbuf.String(), stderrbuf.String(), id, err
}
