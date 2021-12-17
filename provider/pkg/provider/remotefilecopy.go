// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type remotefilecopy struct {
	// Input
	Connection remoteconnection `pulumi:"connection"`
	LocalPath  interface{}      `pulumi:"localPath"`
	RemotePath string           `pulumi:"remotePath"`
}

func (c *remotefilecopy) RunCreate(ctx context.Context, host *provider.HostClient, urn resource.URN) (string, error) {
	host.Log(ctx, diag.Debug, urn,
		fmt.Sprintf("Creating file: %s:%s from local file %s", c.Connection.Host, c.RemotePath, c.LocalPath))

	withConnection := func(f func(*sftp.Client) error) error {
		config, err := c.Connection.SShConfig()
		if err != nil {
			return err
		}
		client, err := c.Connection.Dial(ctx, config)
		if err != nil {
			return err
		}
		defer client.Close()

		sftp, err := sftp.NewClient(client)
		if err != nil {
			return err
		}
		defer sftp.Close()
		return f(sftp)
	}

	inner := func() error {

		switch path := c.LocalPath.(type) {
		case string:
			return withConnection(func(con *sftp.Client) error {
				return c.writeStringPath(con, path)
			})

		case resource.Asset:
			return withConnection(func(con *sftp.Client) error {
				data, err := path.Bytes()
				if err != nil {
					return err
				}
				dstFile, err := con.Create(c.RemotePath)
				if err != nil {
					return err
				}
				reader := bytes.NewReader(data)
				_, err = dstFile.ReadFrom(reader)
				return err
			})

		case resource.Archive:
			return withConnection(func(con *sftp.Client) error {
				reader, err := path.Open()
				if err != nil {
					return err
				}
				defer func() { contract.IgnoreError(reader.Close()) }()

				for src, blob, err := reader.Next(); err != io.EOF; reader.Next() {
					if err != nil {
						return err
					}
					data := make([]byte, blob.Size())

					// TODO: ensure that this complies with the read interface
					read, err := blob.Read(data)
					if err != nil {
						return err
					}
					if read != len(data) {
						return fmt.Errorf("Multipass blob reads are not implemented yet")
					}
					if err := blob.Close(); err != nil {
						return err
					}

					srcPath, name := filepath.Split(src)
					dst := con.Join(append([]string{c.RemotePath}, strings.Split(srcPath, string(filepath.Separator))...)...)

					err = con.MkdirAll(dst)
					if err != nil {
						return err
					}

					dstPath := con.Join(dst, name)
					if err != nil {
						return err
					}

					dstFile, err := con.Create(dstPath)
					if err != nil {
						return err
					}
					reader := bytes.NewReader(data)
					if _, err = dstFile.ReadFrom(reader); err != nil {
						return err
					}
				}
				return nil
			})

		default:
			return fmt.Errorf("Unexpected type: %T", c.LocalPath)
		}
	}
	if err := inner(); err != nil {
		return "", err
	}
	return resource.NewUniqueHex("", 8, 0)
}

func (c *remotefilecopy) writeStringPath(con *sftp.Client, path string) error {
	stats, err := os.Stat(path)
	if err != nil {
		return err
	}

	writeFile := func(src, dst string) error {
		file, err := os.Open(src)
		if err != nil {
			return err
		}
		defer file.Close()

		dstFile, err := con.Create(dst)
		if err != nil {
			return err
		}

		_, err = dstFile.ReadFrom(file)
		return err
	}

	if !stats.IsDir() {
		return writeFile(path, c.RemotePath)
	}

	return filepath.WalkDir(path, func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// We split and then join because we don't know that the remote path has the same
		//separators as the local path.
		dst := con.Join(append([]string{c.RemotePath}, strings.Split(src, string(filepath.Separator))...)...)

		// Files are walked in lexicographic order. That implies that directories come before their
		// contents. That ensures that this is safe.
		if d.IsDir() {
			return con.Mkdir(dst)
		}
		return writeFile(src, dst)

		return nil
	})
}

func (c *remotefilecopy) RunDelete(ctx context.Context, host *provider.HostClient, urn resource.URN) error {
	host.Log(ctx, diag.Debug, urn, "CopyFile delete is a no-op")
	return nil
}
