package remote

import (
	"os"

	"github.com/pkg/sftp"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type CopyFile struct{}

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.CustomResource[CopyFileArgs, CopyFileState])((*CopyFile)(nil))
var _ = (infer.ExplicitDependencies[CopyFileArgs, CopyFileState])((*CopyFile)(nil))

var _ = (infer.Annotated)((*CopyFile)(nil))

// CopyFile implements Annotate which allows you to attach descriptions to the CopyFile resource and its fields.
func (c *CopyFile) Annotate(a infer.Annotator) {
	a.Describe(&c, "Copy a local file to a remote host.")
}

type CopyFileArgs struct {
	Connection *Connection    `pulumi:"connection" provider:"secret"`
	Triggers   *[]interface{} `pulumi:"triggers,optional" providers:"replaceOnDelete"`
	LocalPath  string         `pulumi:"localPath"`
	RemotePath string         `pulumi:"remotePath"`
}

func (c *CopyFileArgs) Annotate(a infer.Annotator) {
	a.Describe(&c.Connection, "The parameters with which to connect to the remote host.")
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.LocalPath, "The path of the file to be copied.")
	a.Describe(&c.RemotePath, "The destination path in the remote host.")
}

func (r *CopyFile) WireDependencies(f infer.FieldSelector, args *CopyFileArgs, state *CopyFileState) {
	f.OutputField(&state.CopyFileArgs.Connection).DependsOn(f.InputField(&args.Connection))
	f.OutputField(&state.CopyFileArgs.Triggers).DependsOn(f.InputField(&args.Triggers))
	f.OutputField(&state.CopyFileArgs.LocalPath).DependsOn(f.InputField(&args.LocalPath))
	f.OutputField(&state.CopyFileArgs.RemotePath).DependsOn(f.InputField(&args.RemotePath))
}

type CopyFileState struct {
	CopyFileArgs
}

func (*CopyFile) Create(ctx p.Context, name string, input CopyFileArgs, preview bool) (string, CopyFileState, error) {
	if preview {
		return "", CopyFileState{input}, nil
	}

	ctx.Logf(diag.Debug,
		"Creating file: %s:%s from local file %s",
		input.Connection.Host, input.RemotePath, input.LocalPath)

	src, err := os.Open(input.LocalPath)
	if err != nil {
		return "", CopyFileState{input}, err
	}
	defer src.Close()

	config, err := input.Connection.SShConfig()
	if err != nil {
		return "", CopyFileState{input}, err
	}
	client, err := input.Connection.Dial(ctx, config)
	if err != nil {
		return "", CopyFileState{input}, err
	}
	defer client.Close()

	sftp, err := sftp.NewClient(client)
	if err != nil {
		return "", CopyFileState{input}, err
	}
	defer sftp.Close()

	dst, err := sftp.Create(input.RemotePath)
	if err != nil {
		return "", CopyFileState{input}, err
	}

	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", CopyFileState{input}, err
	}

	id, err := resource.NewUniqueHex("", 8, 0)
	return id, CopyFileState{input}, err
}
