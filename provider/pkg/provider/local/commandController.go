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
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/gobwas/glob"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/util"
)

// These are not required. They indicate to Go that Command implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors in this file.
var _ = (infer.CustomResource[CommandArgs, CommandState])((*Command)(nil))
var _ = (infer.CustomUpdate[CommandArgs, CommandState])((*Command)(nil))
var _ = (infer.CustomDelete[CommandState])((*Command)(nil))
var _ = (infer.ExplicitDependencies[CommandArgs, CommandState])((*Command)(nil))

func (r *Command) WireDependencies(f infer.FieldSelector, args *CommandArgs, state *CommandState) {

	// get BaseArgs
	interpreterInput := f.InputField(&args.Interpreter)
	dirInput := f.InputField(&args.Dir)
	environmentInput := f.InputField(&args.Environment)
	stdinInput := f.InputField(&args.Stdin)
	assetPathsInput := f.InputField(&args.AssetPaths)
	archivePathsInput := f.InputField(&args.ArchivePaths)

	// set BaseArgs
	f.OutputField(&state.Interpreter).DependsOn(interpreterInput)
	f.OutputField(&state.Dir).DependsOn(dirInput)
	f.OutputField(&state.Environment).DependsOn(environmentInput)
	f.OutputField(&state.Stdin).DependsOn(stdinInput)
	f.OutputField(&state.AssetPaths).DependsOn(assetPathsInput)
	f.OutputField(&state.ArchivePaths).DependsOn(archivePathsInput)

	// set
	triggersInput := f.InputField(&args.Triggers)
	createInput := f.InputField(&args.Create)
	updateInput := f.InputField(&args.Update)
	deleteInput := f.InputField(&args.Delete)

	// set CommandArgs
	f.OutputField(&state.Triggers).DependsOn(triggersInput)
	f.OutputField(&state.Create).DependsOn(createInput)
	f.OutputField(&state.Update).DependsOn(updateInput)
	f.OutputField(&state.Delete).DependsOn(deleteInput)

	f.OutputField(&state.Stdout).DependsOn(
		triggersInput,
		createInput,
		updateInput,
		deleteInput,
		interpreterInput,
		dirInput,
		environmentInput,
		stdinInput,
		assetPathsInput,
		archivePathsInput,
	)
	f.OutputField(&state.Stderr).DependsOn(
		triggersInput,
		createInput,
		updateInput,
		deleteInput,
		interpreterInput,
		dirInput,
		environmentInput,
		stdinInput,
		assetPathsInput,
		archivePathsInput,
	)
}

func (c *Command) Create(ctx p.Context, name string, input CommandArgs, preview bool) (string, CommandState, error) {
	state := CommandState{CommandArgs: input}
	id, err := resource.NewUniqueHex(name, 8, 0)
	if err != nil {
		return id, state, err
	}

	if preview {
		return id, state, nil
	}
	if input.Create == nil {
		return id, state, nil
	}
	cmd := *input.Create
	state.Stdout, state.Stderr, err = state.run(ctx, cmd)
	return id, state, err
}

func (c *Command) Update(ctx p.Context, id string, olds CommandState, news CommandArgs, preview bool) (CommandState, error) {
	state := CommandState{CommandArgs: news}
	if preview {
		return state, nil
	}
	var err error
	if news.Update != nil {
		state.Stdout, state.Stderr, err = state.run(ctx, *news.Update)
	} else if news.Create != nil {
		state.Stdout, state.Stderr, err = state.run(ctx, *news.Create)
	}
	return state, err
}

func (c *Command) Delete(ctx p.Context, id string, props CommandState) error {
	if props.Delete == nil {
		return nil
	}
	_, _, err := props.run(ctx, *props.Delete)
	return err
}

func (c *CommandState) run(ctx p.Context, command string) (string, string, error) {
	var args []string
	if c.Interpreter != nil && len(*c.Interpreter) > 0 {
		args = append(args, *c.Interpreter...)
	} else {
		if runtime.GOOS == "windows" {
			args = []string{"cmd", "/C"}
		} else {
			args = []string{"/bin/sh", "-c"}
		}
	}
	args = append(args, command)

	stdoutr, stdoutw, err := os.Pipe()
	if err != nil {
		return "", "", err
	}
	stderrr, stderrw, err := os.Pipe()
	if err != nil {
		return "", "", err
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = stdoutw
	cmd.Stderr = stderrw
	if c.Dir != nil {
		cmd.Dir = *c.Dir
	} else {
		cmd.Dir, err = os.Getwd()
		if err != nil {
			return "", "", err
		}
	}
	cmd.Env = os.Environ()
	if c.Environment != nil {
		for k, v := range *c.Environment {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	if c.Stdin != nil && len(*c.Stdin) > 0 {
		cmd.Stdin = strings.NewReader(*c.Stdin)
	}

	var stdoutbuf bytes.Buffer
	var stderrbuf bytes.Buffer

	stdouttee := io.TeeReader(stdoutr, &stdoutbuf)
	stderrtee := io.TeeReader(stderrr, &stderrbuf)

	stdoutch := make(chan struct{})
	stderrch := make(chan struct{})
	go util.CopyOutput(ctx, stdouttee, stdoutch, diag.Debug)
	go util.CopyOutput(ctx, stderrtee, stderrch, diag.Error)

	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	stdoutw.Close()
	stderrw.Close()

	<-stdoutch
	<-stderrch

	if err != nil {
		return "", "", err
	}

	if c.AssetPaths != nil {
		assets, err := globAssets(cmd.Dir, *c.AssetPaths)
		if err != nil {
			return "", "", err
		}
		c.Assets = &assets
	}

	if c.ArchivePaths != nil {
		archiveAssets := map[string]interface{}{}
		assets, err := globAssets(cmd.Dir, *c.ArchivePaths)
		if err != nil {
			return "", "", err
		}

		for path, asset := range assets {
			archiveAssets[path] = asset
		}

		archive, err := resource.NewAssetArchive(archiveAssets)
		if err != nil {
			return "", "", err
		}
		c.Archive = archive
	}

	return strings.TrimSuffix(stdoutbuf.String(), "\n"), strings.TrimSuffix(stderrbuf.String(), "\n"), nil
}

func globAssets(dir string, globs []string) (map[string]*resource.Asset, error) {
	assets := map[string]*resource.Asset{}
	compiledGlobs := make([]glob.Glob, len(globs))
	isGlobExclude := make([]bool, len(globs))
	for i, g := range globs {
		isExclude := strings.HasPrefix(g, "!")
		g = strings.TrimPrefix(g, "!")
		compiled, err := glob.Compile(g, '/')
		if err != nil {
			return nil, err
		}
		compiledGlobs[i] = compiled
		isGlobExclude[i] = isExclude
	}

	err := fs.WalkDir(os.DirFS(dir), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		for i, g := range compiledGlobs {
			matched := g.Match(p)
			if !matched {
				continue
			}
			if isGlobExclude[i] {
				delete(assets, p)
			} else {
				assets[p], err = resource.NewPathAsset(path.Join(dir, p))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return assets, nil
}
