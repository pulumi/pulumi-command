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
	"context"
	"fmt"
	"os"

	"github.com/pkg/sftp"

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type remotefilecopy struct {
	// Input
	Connection *remoteconnection `pulumi:"connection"`
	Triggers   *[]interface{}    `pulumi:"triggers,optional"`
	LocalPath  string            `pulumi:"localPath"`
	RemotePath string            `pulumi:"remotePath"`
}

func (c *remotefilecopy) RunCreate(ctx context.Context, host *provider.HostClient, urn resource.URN) (string, error) {

	remoteConnection := c.Connection

	host.Log(ctx, diag.Debug, urn,
		fmt.Sprintf("Creating file: %s:%s from local file %s", remoteConnection.Host, c.RemotePath, c.LocalPath))
	inner := func() error {
		src, err := os.Open(c.LocalPath)
		if err != nil {
			return err
		}
		defer src.Close()

		config, err := remoteConnection.SShConfig()
		if err != nil {
			return err
		}
		client, err := remoteConnection.Dial(ctx, config)
		if err != nil {
			return err
		}
		defer client.Close()

		sftp, err := sftp.NewClient(client)
		if err != nil {
			return err
		}
		defer sftp.Close()

		dst, err := sftp.Create(c.RemotePath)
		if err != nil {
			return err
		}

		_, err = dst.ReadFrom(src)
		return err
	}
	if err := inner(); err != nil {
		return "", err
	}
	return resource.NewUniqueHex("", 8, 0)

}

func (c *remotefilecopy) RunDelete(ctx context.Context, host *provider.HostClient, urn resource.URN) error {
	host.Log(ctx, diag.Debug, urn, "CopyFile delete is a no-op")
	return nil
}
