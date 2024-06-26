// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package remote

import (
	_ "embed"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/infer/types"
)

//go:embed copyToRemote.md
var copyResourceDoc string

type CopyToRemote struct{}

var _ = (infer.Annotated)((*CopyToRemote)(nil))

// Copy implements Annotate which allows you to attach descriptions to the Copy resource.
func (c *CopyToRemote) Annotate(a infer.Annotator) {
	a.Describe(&c, copyResourceDoc)
}

type CopyToRemoteInputs struct {
	Connection *Connection          `pulumi:"connection" provider:"secret"`
	Triggers   *[]interface{}       `pulumi:"triggers,optional" provider:"replaceOnChanges"`
	Source     types.AssetOrArchive `pulumi:"source"`
	RemotePath string               `pulumi:"remotePath"`
}

func (c *CopyToRemoteInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Connection, "The parameters with which to connect to the remote host.")
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.Source, "An [asset or an archive](https://www.pulumi.com/docs/concepts/assets-archives/) "+
		"to upload as the source of the copy. It must be path-based, i.e., be a `FileAsset` or a `FileArchive`. "+
		"The item will be copied as-is; archives like .tgz will not be unpacked. "+
		"Directories are copied recursively, overwriting existing files.")
	a.Describe(&c.RemotePath, "The destination path on the remote host. "+
		"The last element of the path will be created if it doesn't exist but it's an error when additional elements don't exist. "+
		"When the remote path is an existing directory, the source file or directory will be copied into that directory. "+
		"When the source is a file and the remote path is an existing file, that file will be overwritten. "+
		"When the source is a directory and the remote path an existing file, the copy will fail.")
}

func (c *CopyToRemoteInputs) sourcePath() string {
	if c.Source.Asset != nil {
		return c.Source.Asset.Path
	}
	return c.Source.Archive.Path
}

func (c *CopyToRemoteInputs) hash() string {
	if c.Source.Archive != nil {
		return c.Source.Archive.Hash
	}
	return c.Source.Asset.Hash
}

type CopyToRemoteOutputs struct {
	CopyToRemoteInputs
}
