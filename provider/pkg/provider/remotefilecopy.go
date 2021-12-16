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

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type remotefilecopy struct {
	// Input
	Connection remoteconnection `pulumi:"connection"`
	LocalPath  string           `pulumi:"localPath"`
	RemotePath string           `pulumi:"remotePath"`
}

func (c *remotefilecopy) RunCreate(ctx context.Context, host *provider.HostClient, urn resource.URN) (string, error) {
	host.Log(ctx, diag.Info, urn, "Creating remotefilecopy")
	return "", nil
}

func (c *remotefilecopy) RunDelete(ctx context.Context, host *provider.HostClient, urn resource.URN) error {
	host.Log(ctx, diag.Info, urn, "Deleting remotefilecopy")
	return nil
}
