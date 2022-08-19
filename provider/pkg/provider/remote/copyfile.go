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

func (c *CopyFile) Annotate(a infer.Annotator) {
	a.Describe(&c, "Copy a local file to a remote host.")
}

type CopyFileArgs struct {
	Connection Connection     `pulumi:"connection"`
	Triggers   *[]interface{} `pulumi:"triggers,optional"`
	LocalPath  string         `pulumi:"localPath"`
	RemotePath string         `pulumi:"remotePath"`
}

func (c *CopyFileArgs) Annotate(a infer.Annotator) {
	a.Describe(&c.Connection, "The parameters with which to connect to the remote host.")
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.LocalPath, "The path of the file to be copied.")
	a.Describe(&c.RemotePath, "The destination path in the remote host.")
}

type CopyFileState struct{ CopyFileArgs }

func (*CopyFile) Create(ctx p.Context, name string, input CopyFileArgs, preview bool) (string, CopyFileState, error) {
	ctx.Logf(diag.Debug,
		"Creating file: %s:%s from local file %s",
		input.Connection.Host, input.RemotePath, input.LocalPath)
	inner := func() error {
		src, err := os.Open(input.LocalPath)
		if err != nil {
			return err
		}
		defer src.Close()

		config, err := input.Connection.SShConfig()
		if err != nil {
			return err
		}
		client, err := input.Connection.Dial(ctx, config)
		if err != nil {
			return err
		}
		defer client.Close()

		sftp, err := sftp.NewClient(client)
		if err != nil {
			return err
		}
		defer sftp.Close()

		dst, err := sftp.Create(input.RemotePath)
		if err != nil {
			return err
		}

		_, err = dst.ReadFrom(src)
		return err
	}
	if err := inner(); err != nil {
		return "", CopyFileState{input}, err
	}
	id, err := resource.NewUniqueHex("", 8, 0)
	return id, CopyFileState{input}, err
}
