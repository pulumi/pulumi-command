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
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/gobwas/glob"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/util"
)

func run(ctx context.Context, command string, in BaseInputs, out *BaseOutputs, logging *Logging) error {
	contract.Assertf(out != nil, "run:out cannot be nil")
	var args []string
	if in.Interpreter != nil && len(*in.Interpreter) > 0 {
		args = append(args, *in.Interpreter...)
	} else {
		if runtime.GOOS == "windows" {
			args = []string{"cmd", "/C"}
		} else {
			args = []string{"/bin/sh", "-c"}
		}
	}
	args = append(args, command)

	var err error
	var stdoutbuf, stderrbuf, stdouterrbuf bytes.Buffer // stdouterrbuf is only for error messages
	stdouterrwriter := util.ConcurrentWriter{Writer: &stdouterrbuf}
	loggingReader, loggingWriter := io.Pipe()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	stdoutWriters := []io.Writer{&stdoutbuf, &stdouterrwriter}
	if logging.ShouldLogStdout() {
		stdoutWriters = append(stdoutWriters, loggingWriter)
	}
	cmd.Stdout = io.MultiWriter(stdoutWriters...)

	stderrWriters := []io.Writer{&stderrbuf, &stdouterrwriter}
	if logging.ShouldLogStderr() {
		stderrWriters = append(stderrWriters, loggingWriter)
	}
	cmd.Stderr = io.MultiWriter(stderrWriters...)

	if in.Dir != nil {
		cmd.Dir = *in.Dir
	} else {
		cmd.Dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	cmd.Env = os.Environ()
	if in.Environment != nil {
		for k, v := range *in.Environment {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	if in.AddPreviousOutputInEnv == nil || *in.AddPreviousOutputInEnv {
		if out.Stdout != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", util.PULUMI_COMMAND_STDOUT, out.Stdout))
		}
		if out.Stderr != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", util.PULUMI_COMMAND_STDERR, out.Stderr))
		}
	}

	if in.Stdin != nil && len(*in.Stdin) > 0 {
		cmd.Stdin = strings.NewReader(*in.Stdin)
	}

	stdouterrch := make(chan struct{})
	go util.LogOutput(ctx, loggingReader, stdouterrch, diag.Info)

	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	loggingWriter.Close()
	<-stdouterrch

	if err != nil {
		return fmt.Errorf("%w: running %q:\n%s", err, command, stdouterrbuf.String())
	}

	if in.AssetPaths != nil {
		assets, err := globAssets(cmd.Dir, *in.AssetPaths)
		if err != nil {
			return err
		}
		out.Assets = &assets
	}

	if in.ArchivePaths != nil {
		archiveAssets := map[string]interface{}{}
		assets, err := globAssets(cmd.Dir, *in.ArchivePaths)
		if err != nil {
			return err
		}

		for path, asset := range assets {
			archiveAssets[path] = asset
		}

		archive, err := resource.NewAssetArchive(archiveAssets)
		if err != nil {
			return err
		}
		out.Archive = archive
	}

	out.Stdout = strings.TrimSuffix(stdoutbuf.String(), "\n")
	out.Stderr = strings.TrimSuffix(stderrbuf.String(), "\n")

	return nil
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
