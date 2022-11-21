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

package local

import (
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type BaseInputs struct {
	Interpreter  *[]string          `pulumi:"interpreter,optional"`
	Dir          *string            `pulumi:"dir,optional"`
	Environment  *map[string]string `pulumi:"environment,optional"`
	Stdin        *string            `pulumi:"stdin,optional"`
	AssetPaths   *[]string          `pulumi:"assetPaths,optional"`
	ArchivePaths *[]string          `pulumi:"archivePaths,optional"`
}

// Implementing Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *BaseInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Interpreter, "The program and arguments to run the command.\n"+
		"On Linux and macOS, defaults to: `[\"/bin/sh\", \"-c\"]`. On Windows, defaults to: `[\"cmd\", \"/C\"]`")
	a.Describe(&c.Dir, "The directory from which to run the command from. If `dir` does not exist, then\n"+
		"`Command` will fail.")
	a.Describe(&c.Environment, "Additional environment variables available to the command's process.")
	a.Describe(&c.Stdin, "Pass a string to the command's process as standard in")
	a.Describe(&c.ArchivePaths, `A list of path globs to return as a single archive asset after the command completes.

When specifying glob patterns the following rules apply:
- We only include files not directories for assets and archives.
- Path separators are `+"`/`"+` on all platforms - including Windows.
- Patterns starting with `+"`!`"+` are 'exclude' rules.
- Rules are evaluated in order, so exclude rules should be after inclusion rules.
- `+"`*`"+` matches anything except `+"`/`"+`
- `+"`**`"+` matches anything, _including_ `+"`/`"+`
- All returned paths are relative to the working directory (without leading `+"`./`) e.g. `file.text` or `subfolder/file.txt`."+`
- For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)

#### Example

Given the rules:
`+"```yaml"+`
- "assets/**"
- "src/**.js"
- "!**secret.*"
`+"```"+`

When evaluating against this folder:

`+"```yaml"+`
- assets/
  - logos/
    - logo.svg
- src/
  - index.js
  - secret.js
`+"```"+`

The following paths will be returned:

`+"```yaml"+`
- assets/logos/logo.svg
- src/index.js
`+"```")
	a.Describe(&c.AssetPaths, `A list of path globs to read after the command completes.

When specifying glob patterns the following rules apply:
- We only include files not directories for assets and archives.
- Path separators are `+"`/`"+` on all platforms - including Windows.
- Patterns starting with `+"`!`"+` are 'exclude' rules.
- Rules are evaluated in order, so exclude rules should be after inclusion rules.
- `+"`*`"+` matches anything except `+"`/`"+`
- `+"`**`"+` matches anything, _including_ `+"`/`"+`
- All returned paths are relative to the working directory (without leading `+"`./`"+`) e.g. `+"`file.text` or `subfolder/file.txt`"+`.
- For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)

#### Example

Given the rules:
`+"```yaml"+`
- "assets/**"
- "src/**.js"
- "!**secret.*"
`+"```"+`

When evaluating against this folder:

`+"```yaml"+`
- assets/
  - logos/
    - logo.svg
- src/
  - index.js
  - secret.js
`+"```"+`

The following paths will be returned:

`+"```yaml"+`
- assets/logos/logo.svg
- src/index.js
`+"```")
}

type BaseOutputs struct {
	Stdout  string                      `pulumi:"stdout"`
	Stderr  string                      `pulumi:"stderr"`
	Assets  *map[string]*resource.Asset `pulumi:"assets,optional"`
	Archive *resource.Archive           `pulumi:"archive,optional"`
}

// Implementing Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *BaseOutputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Stdout, "The standard output of the command's process")
	a.Describe(&c.Stderr, "The standard error of the command's process")
	a.Describe(&c.Assets, `A map of assets found after running the command.
The key is the relative path from the command dir`)
	a.Describe(&c.Archive, `An archive asset containing files found after running the command.`)
}
