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
	"bufio"
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
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type commandContext struct {
	// Input
	Interpreter  *[]string          `pulumi:"interpreter,optional"`
	Dir          *string            `pulumi:"dir,optional"`
	Environment  *map[string]string `pulumi:"environment,optional"`
	Stdin        *string            `pulumi:"stdin,optional"`
	AssetPaths   *[]string          `pulumi:"assetPaths,optional"`
	ArchivePaths *[]string          `pulumi:"archivePaths,optional"`

	// Output
	Stdout  string                      `pulumi:"stdout"`
	Stderr  string                      `pulumi:"stderr"`
	Assets  *map[string]*resource.Asset `pulumi:"assets,optional"`
	Archive *resource.Archive           `pulumi:"archive,optional"`
}

type run struct {
	commandContext
	Command string `pulumi:"command"`
}

type command struct {
	commandContext
	Triggers *[]interface{} `pulumi:"triggers,optional"`
	Create   string         `pulumi:"create"`
	Delete   *string        `pulumi:"delete,optional"`
	// Optional, if empty will run Create again
	Update *string `pulumi:"update,optional"`
}

func (c *run) RunCommand(ctx context.Context, host *provider.HostClient, urn resource.URN) (string, error) {
	stdout, stderr, id, err := c.run(ctx, c.Command, host, urn)
	c.Stdout = stdout
	c.Stderr = stderr
	return id, err
}

// RunCreate executes the create command, sets Stdout and Stderr, and returns a unique
// ID for the command execution
func (c *command) RunCreate(ctx context.Context, host *provider.HostClient, urn resource.URN) (string, error) {
	stdout, stderr, id, err := c.run(ctx, c.Create, host, urn)
	c.Stdout = stdout
	c.Stderr = stderr
	return id, err
}

//
func (c *command) RunUpdate(ctx context.Context, host *provider.HostClient, urn resource.URN) (string, error) {
	if c.Update != nil {
		stdout, stderr, id, err := c.run(ctx, *c.Update, host, urn)
		c.Stdout = stdout
		c.Stderr = stderr
		return id, err
	}
	stdout, stderr, id, err := c.run(ctx, c.Create, host, urn)
	c.Stdout = stdout
	c.Stderr = stderr
	return id, err
}

// RunDelete executes the create command, sets Stdout and Stderr, and returns a unique
// ID for the command execution
func (c *command) RunDelete(ctx context.Context, host *provider.HostClient, urn resource.URN) error {
	if c.Delete == nil {
		return nil
	}
	_, _, _, err := c.run(ctx, *c.Delete, host, urn)
	return err
}

// run executes the create command, sets Stdout and Stderr, and returns a unique
// ID for the command execution
func (c *commandContext) run(ctx context.Context, command string, host *provider.HostClient, urn resource.URN) (string, string, string, error) {
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
		return "", "", "", err
	}
	stderrr, stderrw, err := os.Pipe()
	if err != nil {
		return "", "", "", err
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = stdoutw
	cmd.Stderr = stderrw
	if c.Dir != nil {
		cmd.Dir = *c.Dir
	} else {
		cmd.Dir, err = os.Getwd()
		if err != nil {
			return "", "", "", err
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
	go copyOutput(ctx, host, urn, stdouttee, stdoutch, diag.Debug)
	go copyOutput(ctx, host, urn, stderrtee, stderrch, diag.Error)

	err = cmd.Start()
	pid := cmd.Process.Pid
	if err == nil {
		err = cmd.Wait()
	}

	stdoutw.Close()
	stderrw.Close()

	<-stdoutch
	<-stderrch

	if err != nil {
		return "", "", "", err
	}

	id, err := resource.NewUniqueHex(fmt.Sprintf("%d", pid), 8, 0)
	if err != nil {
		return "", "", "", err
	}

	if c.AssetPaths != nil {
		assets, err := globAssets(cmd.Dir, *c.AssetPaths)
		if err != nil {
			return "", "", "", err
		}
		c.Assets = &assets
	}

	if c.ArchivePaths != nil {
		archiveAssets := map[string]interface{}{}
		assets, err := globAssets(cmd.Dir, *c.ArchivePaths)
		if err != nil {
			return "", "", "", err
		}

		for path, asset := range assets {
			archiveAssets[path] = asset
		}

		archive, err := resource.NewAssetArchive(archiveAssets)
		if err != nil {
			return "", "", "", err
		}
		c.Archive = archive
	}

	return strings.TrimSuffix(stdoutbuf.String(), "\n"), strings.TrimSuffix(stderrbuf.String(), "\n"), id, nil
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

func copyOutput(ctx context.Context, host *provider.HostClient, urn resource.URN, r io.Reader, doneCh chan<- struct{}, severity diag.Severity) {
	defer close(doneCh)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		err := host.Log(ctx, severity, urn, scanner.Text())
		if err != nil {
			return
		}
	}
}
