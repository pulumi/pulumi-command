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
	"github.com/pulumi/pulumi-go-provider/infer"
)

type CopyFile struct{}

var _ = (infer.Annotated)((*CopyFile)(nil))

// CopyFile implements Annotate which allows you to attach descriptions to the CopyFile resource.
func (c *CopyFile) Annotate(a infer.Annotator) {
	a.Describe(&c, "Copy a local file to a remote host.")
	a.SetResourceDeprecationMessage("This resource is deprecated and will be removed in a future release. " +
		"Please use the `CopyToRemote` resource instead.")
}

type CopyFileInputs struct {
	Connection *Connection    `pulumi:"connection" provider:"secret"`
	Triggers   *[]interface{} `pulumi:"triggers,optional" providers:"replaceOnDelete"`
	LocalPath  string         `pulumi:"localPath"`
	RemotePath string         `pulumi:"remotePath"`
}

// CopyFile implements Annotate which allows you to attach descriptions to the CopyFile resource's fields.
func (c *CopyFileInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Connection, "The parameters with which to connect to the remote host.")
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.LocalPath, "The path of the file to be copied.")
	a.Describe(&c.RemotePath, "The destination path in the remote host.")
}

type CopyFileOutputs struct {
	CopyFileInputs
}
