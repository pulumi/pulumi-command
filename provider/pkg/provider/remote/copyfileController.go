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
	"context"
	"os"

	"github.com/pkg/sftp"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.CustomResource[CopyFileInputs, CopyFileOutputs])((*CopyFile)(nil))

// This is the Create method. This will be run on every CopyFile resource creation.
func (*CopyFile) Create(ctx context.Context, name string, input CopyFileInputs, preview bool) (string, CopyFileOutputs, error) {
	if preview {
		return "", CopyFileOutputs{input}, nil
	}

	p.GetLogger(ctx).Debugf("Creating file: %s:%s from local file %s",
		*input.Connection.Host, input.RemotePath, input.LocalPath)

	src, err := os.Open(input.LocalPath)
	if err != nil {
		return "", CopyFileOutputs{input}, err
	}
	defer src.Close()

	client, err := input.Connection.Dial(ctx)
	if err != nil {
		return "", CopyFileOutputs{input}, err
	}
	defer client.Close()

	sftp, err := sftp.NewClient(client)
	if err != nil {
		return "", CopyFileOutputs{input}, err
	}
	defer sftp.Close()

	dst, err := sftp.Create(input.RemotePath)
	if err != nil {
		return "", CopyFileOutputs{input}, err
	}

	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", CopyFileOutputs{input}, err
	}

	id, err := resource.NewUniqueHex("", 8, 0)
	return id, CopyFileOutputs{input}, err
}
