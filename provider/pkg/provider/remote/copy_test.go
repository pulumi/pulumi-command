// Copyright 2024, Pulumi Corporation.
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
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi-go-provider/infer/types"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/require"
)

func TestAssetSourcePath(t *testing.T) {
	assetPath, input := createAssetInput(t)
	require.NotNil(t, input.Source.Asset)
	require.Equal(t, assetPath, input.Source.Asset.Path)
	require.Equal(t, assetPath, input.sourcePath())
}

func TestArchiveSourcePath(t *testing.T) {
	archivePath, input := createArchiveInput(t)
	require.NotNil(t, input.Source.Archive)
	require.Equal(t, archivePath, input.Source.Archive.Path)
	require.Equal(t, archivePath, input.sourcePath())
}

func TestAssetHash(t *testing.T) {
	_, input := createAssetInput(t)
	require.NotNil(t, input.Source.Asset)
	require.Equal(t, input.Source.Asset.Hash, input.hash())
}

func TestArchiveHash(t *testing.T) {
	_, input := createArchiveInput(t)
	require.NotNil(t, input.Source.Archive)
	require.Equal(t, input.Source.Archive.Hash, input.hash())
}

func createArchiveInput(t *testing.T) (string, *CopyToRemoteInputs) {
	archivePath := filepath.Join(t.TempDir(), "archive.zip")
	require.NoError(t, os.WriteFile(archivePath, []byte("hello, world"), 0644))
	archive, err := resource.NewPathArchive(archivePath)
	require.NoError(t, err)

	c := &CopyToRemoteInputs{
		Source: types.AssetOrArchive{
			Archive: archive,
		},
	}
	return archivePath, c
}

func createAssetInput(t *testing.T) (string, *CopyToRemoteInputs) {
	assetPath := filepath.Join(t.TempDir(), "asset")
	require.NoError(t, os.WriteFile(assetPath, []byte("hello, world"), 0644))
	asset, err := resource.NewPathAsset(assetPath)
	require.NoError(t, err)

	c := &CopyToRemoteInputs{
		Source: types.AssetOrArchive{
			Asset: asset,
		},
	}
	return assetPath, c
}
