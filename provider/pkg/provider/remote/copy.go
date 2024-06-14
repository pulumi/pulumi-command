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
	"github.com/pulumi/pulumi-go-provider/infer/types"
)

type Copy struct{}

var _ = (infer.Annotated)((*Copy)(nil))

// Copy implements Annotate which allows you to attach descriptions to the Copy resource.
func (c *Copy) Annotate(a infer.Annotator) {
	a.Describe(&c, "Copy an Asset or Archive to a remote host.")
}

type CopyInputs struct {
	Connection *Connection          `pulumi:"connection" provider:"secret"`
	Triggers   *[]interface{}       `pulumi:"triggers,optional" provider:"replaceOnChanges"`
	Source     types.AssetOrArchive `pulumi:"source"`
	RemotePath string               `pulumi:"remotePath"`
}

func (c *CopyInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Connection, "The parameters with which to connect to the remote host.")
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.Source, "An asset or an archive to upload as the source of the copy. It must be path based.")
	a.Describe(&c.RemotePath, "The destination path in the remote host.")
}

func (c *CopyInputs) sourcePath() string {
	if c.Source.Asset != nil {
		return c.Source.Asset.Path
	}
	return c.Source.Archive.Path
}

func (c *CopyInputs) hash() string {
	if c.Source.Archive != nil {
		return c.Source.Archive.Hash
	}
	return c.Source.Asset.Hash
}

type CopyOutputs struct {
	CopyInputs
}
