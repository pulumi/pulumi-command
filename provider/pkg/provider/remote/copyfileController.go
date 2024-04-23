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
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.CustomResource[CopyFileInputs, CopyFileOutputs])((*CopyFile)(nil))

// This is the Create method. This will be run on every CopyFile resource creation.
func (*CopyFile) Create(ctx p.Context, name string, input CopyFileInputs, preview bool) (string, CopyFileOutputs, error) {
	if preview {
		return "", CopyFileOutputs{input}, nil
	}

	ctx.Logf(diag.Debug,
		"Creating: %s:%s from local '%s'",
		input.Connection.Host, input.RemotePath, input.LocalPath)

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

	src, err := os.Open(input.LocalPath)
	if err != nil {
		return "", CopyFileOutputs{input}, err
	}
	defer src.Close()

	srcInfo, err := src.Stat()
	if err != nil {
		return "", CopyFileOutputs{input}, err
	}
	if srcInfo.IsDir() {
		err = copyDir(sftp, input.LocalPath, input.RemotePath)
	} else {
		err = copyFile(sftp, input.LocalPath, input.RemotePath)
	}
	if err != nil {
		return "", CopyFileOutputs{input}, err
	}

	id, err := resource.NewUniqueHex("", 8, 0)
	return id, CopyFileOutputs{input}, err
}

func copyFile(sftp *sftp.Client, src, dst string) error {
	local, err := os.Open(src)
	if err != nil {
		return err
	}
	defer local.Close()

	remote, err := sftp.Create(dst)
	if err != nil {
		return err
	}
	defer remote.Close()

	_, err = remote.ReadFrom(local)
	return err
}

func copyDir(sftp *sftp.Client, src, dst string) error {
	fileSystem := os.DirFS(src)
	return fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		remotePath := filepath.Join(dst, path)

		if d.IsDir() {
			return sftp.Mkdir(remotePath)
		}
		return copyFile(sftp, filepath.Join(src, path), remotePath)
	})
}
