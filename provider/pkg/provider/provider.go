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
	"strings"

	"github.com/blang/semver"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/local"
	"github.com/pulumi/pulumi-command/provider/pkg/provider/remote"
	"github.com/pulumi/pulumi-go-provider/integration"
)

func Provider() p.Provider {
	return infer.NewProvider().
		WithResources(
			infer.Resource[*local.Command, local.CommandArgs, local.CommandState](),
			infer.Resource[*remote.Command, remote.CommandArgs, remote.CommandState](),
			infer.Resource[*remote.CopyFile, remote.CopyFileArgs, remote.CopyFileState](),
		).
		WithFunctions(
			infer.Function[*local.Run, local.RunArgs, local.RunState](),
		).
		WithDisplayName("Command").
		WithDescription("The Pulumi Command Provider enables you to execute commands and scripts either locally or remotely as part of the Pulumi resource model.").
		WithKeywords([]string{
			"pulumi",
			"command",
			"category/utility",
			"kind/native"}).
		WithHomepage("https://pulumi.com").
		WithLicense("Apache-2.0").
		WithRepository("https://github.com/pulumi/pulumi-command").
		WithPublisher("Pulumi").
		WithLogoURL("https://raw.githubusercontent.com/pulumi/pulumi-command/master/assets/logo.svg").
		WithLanguageMap(map[string]any{
			"csharp": map[string]any{
				"packageReferences": map[string]string{
					"Pulumi": "3.*",
				},
			},
			"go": map[string]any{
				"generateResourceContainerTypes": true,
				"importBasePath":                 "github.com/pulumi/pulumi-command/sdk/go/command",
			},
			"nodejs": map[string]any{
				"dependencies": map[string]string{
					"@pulumi/pulumi": "^3.0.0",
				},
			},
			"python": map[string]any{
				"requires": map[string]string{
					"pulumi": ">=3.0.0,<4.0.0",
				},
			},
			"java": map[string]any{
				"buildFiles":                      "gradle",
				"gradleNexusPublishPluginVersion": "1.1.0",
				"dependencies": map[string]any{
					"com.pulumi:pulumi":               "0.6.0",
					"com.google.code.gson:gson":       "2.8.9",
					"com.google.code.findbugs:jsr305": "3.0.2",
				},
			},
		})
}

func Schema(version string) (string, error) {
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}
	s, err := integration.NewServer("command", semver.MustParse(version), Provider()).
		GetSchema(p.GetSchemaRequest{})
	return s.Schema, err
}
