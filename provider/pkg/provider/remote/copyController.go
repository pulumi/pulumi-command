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
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.CustomResource[CopyInputs, CopyOutputs])((*Copy)(nil))

func doCopy(ctx context.Context, input CopyInputs) (CopyOutputs, error) {
	sourcePath := input.sourcePath()

	p.GetLogger(ctx).Debugf("Creating file: %s:%s from local file %s",
		*input.Connection.Host, input.RemotePath, sourcePath)

	src, err := os.Open(sourcePath)
	if err != nil {
		return CopyOutputs{input}, err
	}
	defer src.Close()

	client, err := input.Connection.Dial(ctx)
	if err != nil {
		return CopyOutputs{input}, err
	}
	defer client.Close()

	sftp, err := sftp.NewClient(client)
	if err != nil {
		return CopyOutputs{input}, err
	}
	defer sftp.Close()

	srcInfo, err := src.Stat()
	if err != nil {
		return CopyOutputs{input}, err
	}
	if srcInfo.IsDir() {
		err = copyDir(sftp, sourcePath, input.RemotePath)
	} else {
		err = copyFile(sftp, sourcePath, input.RemotePath)
	}
	return CopyOutputs{input}, err
}

// This is the Create method. This will be run on every Copy resource creation.
func (*Copy) Create(ctx context.Context, name string, input CopyInputs, preview bool) (string, CopyOutputs, error) {
	if preview {
		return "", CopyOutputs{input}, nil
	}

	outputs, err := doCopy(ctx, input)
	if err != nil {
		return "", CopyOutputs{input}, err
	}

	id, err := resource.NewUniqueHex("", 8, 0)
	return id, outputs, err
}

func (c *Copy) Update(ctx context.Context, id string, olds CopyOutputs, news CopyInputs, preview bool) (CopyOutputs, error) {
	if err := news.validate(); err != nil {
		return olds, err
	}

	needCopy := true
	if news.hash() == olds.hash() && news.RemotePath == olds.RemotePath {
		needCopy = false
	}

	if preview {
		// TODO how to show the user whether we need to copy or not?
		return CopyOutputs{news}, nil
	}

	if needCopy {
		return doCopy(ctx, news)
	}
	return CopyOutputs{news}, nil
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

// copyDir copies a directory recursively from the local file system to a remote host.
// Note that the current is naive and sequential and therefore can be slow.
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
