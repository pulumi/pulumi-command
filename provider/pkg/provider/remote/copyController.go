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
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/archive"
)

// copyTextContent writes text content directly to a remote file via SFTP.
func copyTextContent(sftpClient *sftp.Client, content, destPath string) error {
	destStat, err := remoteStat(sftpClient, destPath)
	if err != nil {
		return err
	}

	// If destination is a directory, we cannot copy text content to it without a filename.
	// The user must provide a full file path as the destination.
	if destStat != nil && destStat.IsDir() {
		return fmt.Errorf("remote path %s is a directory; when using a text asset, remotePath must be a file path", destPath)
	}

	// Ensure parent directories exist.
	if err := sftpClient.MkdirAll(filepath.Dir(destPath)); err != nil {
		return fmt.Errorf("failed to create parent directories for %s: %w", destPath, err)
	}

	remote, err := sftpClient.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create remote file %s: %w", destPath, err)
	}
	defer remote.Close()

	_, err = remote.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to write text content to remote path %s: %w", destPath, err)
	}
	return nil
}

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var (
	_ = (infer.CustomResource[CopyToRemoteInputs, CopyToRemoteOutputs])((*CopyToRemote)(nil))
	_ = (infer.CustomCheck[CopyToRemoteInputs])((*CopyToRemote)(nil))
	_ = (infer.CustomUpdate[CopyToRemoteInputs, CopyToRemoteOutputs])((*CopyToRemote)(nil))
)

func (c *CopyToRemote) Check(
	ctx context.Context,
	req infer.CheckRequest,
) (infer.CheckResponse[CopyToRemoteInputs], error) {
	var failures []p.CheckFailure

	newInputs := req.NewInputs
	inputs, newFailures, err := infer.DefaultCheck[CopyToRemoteInputs](ctx, newInputs)
	failures = append(failures, newFailures...)
	if err != nil {
		return infer.CheckResponse[CopyToRemoteInputs]{Inputs: inputs, Failures: failures}, err
	}

	// If source is unknown (computed during preview), skip asset/archive validation
	// since the value isn't available yet.
	sourceVal, sourceOk := newInputs.GetOk("source")
	sourceIsComputed := sourceOk && sourceVal.IsComputed()

	if !sourceIsComputed {
		hasAsset := inputs.Source.Asset != nil
		hasArchive := inputs.Source.Archive != nil

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
	}

	return infer.CheckResponse[CopyToRemoteInputs]{Inputs: inputs, Failures: failures}, nil
}

// This is the Create method. This will be run on every Copy resource creation.
func (*CopyToRemote) Create(
	ctx context.Context,
	req infer.CreateRequest[CopyToRemoteInputs],
) (infer.CreateResponse[CopyToRemoteOutputs], error) {
	input := req.Inputs
	preview := req.DryRun
	if preview {
		return infer.CreateResponse[CopyToRemoteOutputs]{ID: "", Output: CopyToRemoteOutputs{input}}, nil
	}

	outputs, err := copyToRemote(ctx, input)
	if err != nil {
		return infer.CreateResponse[CopyToRemoteOutputs]{ID: "", Output: CopyToRemoteOutputs{input}}, err
	}

	id, err := resource.NewUniqueHex("", 8, 0)
	return infer.CreateResponse[CopyToRemoteOutputs]{ID: id, Output: outputs}, err
}

func (c *CopyToRemote) Update(
	ctx context.Context,
	req infer.UpdateRequest[CopyToRemoteInputs, CopyToRemoteOutputs],
) (infer.UpdateResponse[CopyToRemoteOutputs], error) {
	olds := req.State
	news := req.Inputs
	preview := req.DryRun
	if preview {
		return infer.UpdateResponse[CopyToRemoteOutputs]{Output: CopyToRemoteOutputs{news}}, nil
	}

	needCopy := news.hash() != olds.hash() || news.RemotePath != olds.RemotePath
	if needCopy {
		outputs, err := copyToRemote(ctx, news)
		return infer.UpdateResponse[CopyToRemoteOutputs]{Output: outputs}, err
	}
	return infer.UpdateResponse[CopyToRemoteOutputs]{Output: CopyToRemoteOutputs{news}}, nil
}

// copyToRemote unpacks the inputs, dials the SSH connection, creates an sFTP client, and dispatches
// to the appropriate copy routine based on the source asset/archive subtype.
func copyToRemote(ctx context.Context, input CopyToRemoteInputs) (CopyToRemoteOutputs, error) {
	p.GetLogger(ctx).Debugf("Creating %s:%s from %s",
		*input.Connection.Host, input.RemotePath, sourceDescription(input))

	client, err := input.Connection.Dial(ctx)
	if err != nil {
		return CopyToRemoteOutputs{input}, err
	}
	defer client.Close()

	// The docs warns that concurrent writes "require special consideration. A write to a later
	/// offset in a file after an error, could end up with a file length longer than what was
	// successfully written."
	// We don't do subsequent writes to the same file, only a single ReadFrom, so we should be fine.
	sftpClient, err := sftp.NewClient(client, sftp.UseConcurrentWrites(true))
	if err != nil {
		return CopyToRemoteOutputs{input}, err
	}
	defer sftpClient.Close()

	if input.Source.Asset != nil {
		err = copyAssetToRemote(sftpClient, input.Source.Asset, input.RemotePath)
	} else {
		err = copyArchiveToRemote(sftpClient, input.Source.Archive, input.RemotePath)
	}
	return CopyToRemoteOutputs{input}, err
}

func sourceDescription(input CopyToRemoteInputs) string {
	if a := input.Source.Asset; a != nil {
		switch {
		case a.IsText():
			return "text asset"
		case a.IsPath():
			return fmt.Sprintf("local file %s", a.Path)
		case a.IsURI():
			return fmt.Sprintf("remote asset %s", a.URI)
		}
	}
	if a := input.Source.Archive; a != nil {
		switch {
		case a.IsPath():
			return fmt.Sprintf("local path %s", a.Path)
		case a.IsURI():
			return fmt.Sprintf("remote archive %s", a.URI)
		case a.IsAssets():
			return "asset archive"
		}
	}
	return "unknown source"
}

func copyAssetToRemote(sftpClient *sftp.Client, a *resource.Asset, destPath string) error {
	switch {
	case a.IsText():
		return copyTextContent(sftpClient, a.Text, destPath)
	case a.IsPath():
		return sftpCopy(sftpClient, a.Path, destPath)
	case a.IsURI():
		blob, err := a.Read()
		if err != nil {
			return fmt.Errorf("failed to read remote asset %s: %w", a.URI, err)
		}
		defer blob.Close()
		return copyReaderAsFile(sftpClient, blob, uriBasename(a.URI), destPath)
	}
	return fmt.Errorf("asset is neither path-based, text-based, nor URI-based")
}

func copyArchiveToRemote(sftpClient *sftp.Client, a *resource.Archive, destPath string) error {
	switch {
	case a.IsPath():
		return sftpCopy(sftpClient, a.Path, destPath)
	case a.IsURI():
		format, rc, err := a.ReadSourceArchive()
		if err != nil {
			return fmt.Errorf("failed to read remote archive %s: %w", a.URI, err)
		}
		if format == archive.NotArchive || rc == nil {
			return fmt.Errorf("URL %q is not a recognized archive format", a.URI)
		}
		defer rc.Close()
		return copyReaderAsFile(sftpClient, rc, uriBasename(a.URI), destPath)
	case a.IsAssets():
		return copyAssetArchive(sftpClient, a, destPath)
	}
	return fmt.Errorf("archive is neither path-based, URI-based, nor asset-based")
}

// copyReaderAsFile streams r to destPath via SFTP, mirroring the file-handling logic of sftpCopy:
// when destPath is an existing directory the contents are written to destPath/sourceName, when
// destPath does not exist the parent directories are created and the file is written at destPath,
// and when destPath is an existing file it is overwritten.
func copyReaderAsFile(sftpClient *sftp.Client, r io.Reader, sourceName, destPath string) error {
	destStat, err := remoteStat(sftpClient, destPath)
	if err != nil {
		return err
	}

	dest := destPath
	if destStat != nil && destStat.IsDir() {
		if sourceName == "" {
			return fmt.Errorf("remote path %s is a directory; cannot determine destination filename", destPath)
		}
		dest = filepath.Join(dest, sourceName)
	} else if destStat == nil {
		if err := sftpClient.MkdirAll(filepath.Dir(dest)); err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %w", dest, err)
		}
	}

	remote, err := sftpClient.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create remote file %s: %w", dest, err)
	}
	defer remote.Close()

	if _, err := remote.ReadFrom(r); err != nil {
		return fmt.Errorf("failed to write to remote path %s: %w", dest, err)
	}
	return nil
}

// copyAssetArchive iterates over the entries of an AssetArchive and writes each one to destPath/name.
func copyAssetArchive(sftpClient *sftp.Client, a *resource.Archive, destPath string) error {
	destStat, err := remoteStat(sftpClient, destPath)
	if err != nil {
		return err
	}
	if destStat != nil && !destStat.IsDir() {
		return fmt.Errorf("remote path %s exists but is not a directory; cannot copy asset archive contents", destPath)
	}
	if destStat == nil {
		if err := sftpClient.MkdirAll(destPath); err != nil {
			return fmt.Errorf("failed to create remote directory %s: %w", destPath, err)
		}
	}

	reader, err := a.Open()
	if err != nil {
		return fmt.Errorf("failed to open asset archive: %w", err)
	}
	defer reader.Close()

	for {
		name, blob, err := reader.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read asset archive entry: %w", err)
		}
		if err := writeArchiveEntry(sftpClient, blob, filepath.Join(destPath, name)); err != nil {
			return err
		}
	}
}

func writeArchiveEntry(sftpClient *sftp.Client, blob io.ReadCloser, remotePath string) error {
	defer blob.Close()
	if err := sftpClient.MkdirAll(filepath.Dir(remotePath)); err != nil {
		return fmt.Errorf("failed to create parent directories for %s: %w", remotePath, err)
	}
	remote, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file %s: %w", remotePath, err)
	}
	defer remote.Close()
	if _, err := remote.ReadFrom(blob); err != nil {
		return fmt.Errorf("failed to copy archive entry to %s: %w", remotePath, err)
	}
	return nil
}

// uriBasename returns the basename of the path component of the given URI, or "" when one cannot be
// determined (e.g. the URI is malformed or has no path).
func uriBasename(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	base := path.Base(u.Path)
	if base == "/" || base == "." {
		return ""
	}
	return base
}

// If the file does not exist, returns nil, nil.
func remoteStat(sftpClient *sftp.Client, path string) (fs.FileInfo, error) {
	info, err := sftpClient.Stat(path)
	// sftp normalizes the error to os.ErrNotExist, see client.go: normaliseError.
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to stat remote path %s: %w", path, err)
	}
	return info, nil
}

func sftpCopy(sftpClient *sftp.Client, sourcePath, destPath string) error {
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
	destStat, err = remoteStat(sftpClient, destPath)
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
			err = sftpClient.MkdirAll(dest)
			if err != nil {
				return fmt.Errorf("failed to create remote directory %s: %w", dest, err)
			}
		}

		if !strings.HasSuffix(sourcePath, "/") {
			dest = filepath.Join(dest, filepath.Base(sourcePath))
			destStat, err := remoteStat(sftpClient, dest)
			if err != nil {
				return err
			}
			// It's ok if the dir exists, we'll copy into it.
			if destStat != nil && !destStat.IsDir() {
				return fmt.Errorf("remote path %s exists but is not a directory", dest)
			}
			if destStat == nil {
				err = sftpClient.Mkdir(dest)
				if err != nil {
					return fmt.Errorf("failed to create remote directory %s: %w", dest, err)
				}
			}
		}
		err = copyDir(sftpClient, sourcePath, dest)
	} else {
		// If the file is f and the destination is existing dir/, copy to dir/f.
		if destStat != nil && destStat.IsDir() {
			dest = filepath.Join(dest, filepath.Base(sourcePath))
		} else if destStat == nil {
			// Ensure parent directories exist when creating a new file.
			if err := sftpClient.MkdirAll(filepath.Dir(dest)); err != nil {
				return fmt.Errorf("failed to create parent directories for %s: %w", dest, err)
			}
		}
		err = copyFile(sftpClient, sourcePath, dest)
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

		dirInfo, err := remoteStat(sftp, remotePath)
		if err != nil {
			return err
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
