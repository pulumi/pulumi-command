package remote

import (
	"github.com/pulumi/pulumi-go-provider/infer"
)

type CopyFile struct{}

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

type CopyFileState struct {
	CopyFileArgs
}
