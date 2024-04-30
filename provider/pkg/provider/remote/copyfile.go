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
	"errors"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/archive"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/asset"
)

type CopyFile struct{}

var _ = (infer.Annotated)((*CopyFile)(nil))

// Copy implements Annotate which allows you to attach descriptions to the Copy resource.
func (c *CopyFile) Annotate(a infer.Annotator) {
	a.Describe(&c, "Copy an Asset or Archive to a remote host.")
}

type CopyFileInputs struct {
	Connection   *Connection      `pulumi:"connection" provider:"secret"`
	Triggers     *[]interface{}   `pulumi:"triggers,optional" providers:"replaceOnDelete"`
	LocalAsset   *asset.Asset     `pulumi:"localAsset,optional"`
	LocalArchive *archive.Archive `pulumi:"localArchive,optional"`
	RemotePath   string           `pulumi:"remotePath"`
}

// CopyFile implements Annotate which allows you to attach descriptions to the CopyFile resource's fields.
func (c *CopyFileInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Connection, "The parameters with which to connect to the remote host.")
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.LocalAsset, "The path of the file to be copied. Only one of LocalAsset or LocalArchive can be set.")
	a.Describe(&c.LocalArchive, "The path of the folder or archive to be copied. Only one of LocalAsset or LocalArchive can be set.")
	a.Describe(&c.RemotePath, "The destination path in the remote host.")
}

func (c *CopyFileInputs) validate() error {
	if c.LocalAsset != nil && c.LocalArchive != nil {
		return errors.New("only one of LocalAsset or LocalArchive can be set")
	}
	if c.LocalAsset == nil && c.LocalArchive == nil {
		return errors.New("either LocalAsset or LocalArchive must be set")
	}
	if c.LocalAsset != nil && !c.LocalAsset.IsPath() {
		return errors.New("LocalAsset must be a file asset")
	}
	if c.LocalArchive != nil && !c.LocalArchive.IsPath() {
		return errors.New("LocalArchive must be a path to a file or directory")
	}
	return nil
}

func (c *CopyFileInputs) sourcePath() string {
	if c.LocalAsset != nil {
		return c.LocalAsset.Path
	}
	return c.LocalArchive.Path
}

type CopyFileOutputs struct {
	CopyFileInputs
}
