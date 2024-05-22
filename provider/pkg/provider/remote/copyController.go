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
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.CustomResource[CopyInputs, CopyOutputs])((*Copy)(nil))
var _ = (infer.CustomCheck[CopyInputs])((*Copy)(nil))
var _ = (infer.CustomUpdate[CopyInputs, CopyOutputs])((*Copy)(nil))

func (c *Copy) Check(ctx context.Context, urn string, oldInputs, newInputs resource.PropertyMap) (CopyInputs, []p.CheckFailure, error) {
	var failures []p.CheckFailure

	hasAsset := newInputs.HasValue("asset")
	hasArchive := newInputs.HasValue("archive")

	if hasAsset && hasArchive {
		failures = append(failures, p.CheckFailure{
			Property: "asset",
			Reason:   "only one of asset or archive can be set",
		})
	}
	if !hasAsset && !hasArchive {
		failures = append(failures, p.CheckFailure{
			Property: "asset",
			Reason:   "either asset or archive must be set",
		})
	}

	inputs, newFailures, err := infer.DefaultCheck[CopyInputs](newInputs)
	if err != nil {
		return inputs, failures, err
	}
	failures = append(failures, newFailures...)

	if hasAsset && !inputs.Asset.IsPath() {
		failures = append(failures, p.CheckFailure{
			Property: "asset",
			Reason:   "asset must be a path-based file asset",
		})
	}
	if hasArchive && !inputs.Archive.IsPath() {
		failures = append(failures, p.CheckFailure{
			Property: "archive",
			Reason:   "archive must be a path to a file or directory",
		})
	}

	return inputs, failures, nil
}

// This is the Create method. This will be run on every Copy resource creation.
func (*Copy) Create(ctx context.Context, name string, input CopyInputs, preview bool) (string, CopyOutputs, error) {
	if preview {
		return "", CopyOutputs{input}, nil
	}

	outputs, err := copy(ctx, input)
	if err != nil {
		return "", CopyOutputs{input}, err
	}

	id, err := resource.NewUniqueHex("", 8, 0)
	return id, outputs, err
}

func (c *Copy) Update(ctx context.Context, id string, olds CopyOutputs, news CopyInputs, preview bool) (CopyOutputs, error) {
	if preview {
		return CopyOutputs{news}, nil
	}

	needCopy := news.hash() != olds.hash() || news.RemotePath != olds.RemotePath
	if needCopy {
		return copy(ctx, news)
	}
	return CopyOutputs{news}, nil
}

// copy unpacks the inputs, dials the SSH connection, creates an sFTP client, and calls sftpCopy.
func copy(ctx context.Context, input CopyInputs) (CopyOutputs, error) {
	sourcePath := input.sourcePath()

	p.GetLogger(ctx).Debugf("Creating file: %s:%s from local file %s",
		*input.Connection.Host, input.RemotePath, sourcePath)

	client, err := input.Connection.Dial(ctx)
	if err != nil {
		return CopyOutputs{input}, err
	}
	defer client.Close()

	// The docs warns that concurrent writes "require special consideration. A write to a later
	/// offset in a file after an error, could end up with a file length longer than what was
	// successfully written."
	// We don't do subsequent writes to the same file, only a single ReadFrom, so we should be fine.
	sftp, err := sftp.NewClient(client, sftp.UseConcurrentWrites(true))
	if err != nil {
		return CopyOutputs{input}, err
	}
	defer sftp.Close()

	err = sftpCopy(sftp, sourcePath, input.RemotePath)
	return CopyOutputs{input}, err
}

// If the file does not exist, returns nil, nil.
func remoteStat(sftp *sftp.Client, path string) (fs.FileInfo, error) {
	info, err := sftp.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to stat remote path %s: %w", path, err)
	}
	return info, nil
}

func sftpCopy(sftp *sftp.Client, sourcePath, destPath string) error {
	src, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	srcInfo, err := src.Stat()
	if err != nil {
		return err
	}

	var destStat fs.FileInfo
	destStat, err = remoteStat(sftp, destPath)
	if err != nil {
		return err
	}

	// Before copying, we might need to adjust some paths. Files have different semantics from
	// directories, and source directories depend on whether they have a trailing slash.
	//
	// source | dest - exists as dir | dest - does not exist | dest - exists as file
	// -------|----------------------|-----------------------|-----------------------
	// dir    | dest/dir             | dest/dir              | error
	// dir/   | dest/x for x in dir  | dest/dir              | error
	// file   | dest/file            | dest                  | dest (overwritten)
	dest := destPath
	if srcInfo.IsDir() {
		if destStat == nil {
			err = sftp.Mkdir(dest)
			if err != nil {
				return fmt.Errorf("failed to create remote directory %s: %w", dest, err)
			}
		}

		if !strings.HasSuffix(sourcePath, "/") {
			dest = filepath.Join(dest, filepath.Base(sourcePath))
			err = sftp.Mkdir(dest)
			if err != nil {
				return fmt.Errorf("failed to create remote directory %s: %w", dest, err)
			}
		}
		err = copyDir(sftp, sourcePath, dest)
	} else {
		// If the file is f and the destination is existing dir/, copy to dir/f.
		if destStat != nil && destStat.IsDir() {
			dest = filepath.Join(dest, filepath.Base(sourcePath))
		}
		err = copyFile(sftp, sourcePath, dest)
	}
	return err
}

func copyFile(sftp *sftp.Client, src, dst string) error {
	local, err := os.Open(src)
	if err != nil {
		return err
	}
	defer local.Close()

	remote, err := sftp.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create remote file %s: %w", dst, err)
	}
	defer remote.Close()

	_, err = remote.ReadFrom(local)
	if err != nil {
		return fmt.Errorf("failed to copy file %s to remote path %s: %w", src, dst, err)
	}
	return nil
}

// copyDir copies a directory recursively from the local file system to a remote host.
// Note that the current implementation is naive and sequential and therefore can be slow.
func copyDir(sftp *sftp.Client, src, dst string) error {
	fileSystem := os.DirFS(src)
	return fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		remotePath := filepath.Join(dst, path)

		if !d.IsDir() {
			return copyFile(sftp, filepath.Join(src, path), remotePath)
		}

		dirInfo, err := sftp.Stat(remotePath)
		// sftp normalizes the error to os.ErrNotExist, see client.go: normaliseError.
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to stat remote path %s: %w", remotePath, err)
		}

		if dirInfo == nil {
			if err = sftp.Mkdir(remotePath); err != nil {
				return fmt.Errorf("failed to create remote directory %s: %w", remotePath, err)
			}
		} else if !dirInfo.IsDir() {
			return fmt.Errorf("remote path %s exists but is not a directory", remotePath)
		}
		return nil
	})
}
