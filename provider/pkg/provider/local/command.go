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

type Command struct{}

var _ = (infer.CustomResource[CommandArgs, CommandState])((*Command)(nil))
var _ = (infer.CustomUpdate[CommandArgs, CommandState])((*Command)(nil))
var _ = (infer.CustomDelete[CommandState])((*Command)(nil))

func (c *Command) Annotate(a infer.Annotator) {
	a.Describe(&c, `A local command to be executed.
This command can be inserted into the life cycles of other resources using the
`+"`dependsOn`"+` or `+"`parent`"+` resource options. A command is considered to have
failed when it finished with a non-zero exit code. This will fail the CRUD step
of the `+"`Command`"+` resource.`)
}

type BaseArgs struct {
	Interpreter  *[]string               `pulumi:"interpreter,optional"`
	Dir          *string                 `pulumi:"dir,optional"`
	Environment  *map[string]interface{} `pulumi:"environment,optional"`
	Stdin        *string                 `pulumi:"stdin,optional"`
	AssetPaths   *[]string               `pulumi:"assetPaths,optional"`
	ArchivePaths *[]string               `pulumi:"archivePaths,optional"`
}

func (c *BaseArgs) Annotate(a infer.Annotator) {
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

type BaseState struct {
	Stdout  string                      `pulumi:"stdout"`
	Stderr  string                      `pulumi:"stderr"`
	Assets  *map[string]*resource.Asset `pulumi:"assets,optional"`
	Archive *resource.Archive           `pulumi:"archive,optional"`
}

func (c *BaseState) Annotate(a infer.Annotator) {
	a.Describe(&c.Stdout, "The standard output of the command's process")
	a.Describe(&c.Stderr, "The standard error of the command's process")
	a.Describe(&c.Assets, `A map of assets found after running the command.
The key is the relative path from the command dir`)
	a.Describe(&c.Archive, `An archive asset containing files found after running the command.`)
}

type CommandArgs struct {
	BaseArgs
	Triggers *[]any  `pulumi:"triggers,optional" provider:"replaceOnChanges"`
	Create   *string `pulumi:"create,optional"`
	Delete   *string `pulumi:"delete,optional"`
	Update   *string `pulumi:"update,optional"`
}

func (c *CommandArgs) Annotate(a infer.Annotator) {
	a.Describe(&c.Triggers, "Trigger replacements on changes to this input.")
	a.Describe(&c.Create, "The command to run on create.")
	a.Describe(&c.Delete, "The command to run on delete.")
	a.Describe(&c.Update, "The command to run on update, if empty, create will run again.")
}

type CommandState struct {
	CommandArgs
	BaseState
}

func (c *Command) Create(ctx p.Context, name string, input CommandArgs, preview bool) (id string, output CommandState, err error) {
	state := CommandState{CommandArgs: input}
	if preview {
		return "", state, nil
	}
	cmd := ""
	if input.Create != nil {
		cmd = *input.Create
	}
	stdout, stderr, id, err := state.run(ctx, cmd)
	state.Stdout = stdout
	state.Stderr = stderr
	return id, state, err
}

func (c *Command) Update(ctx p.Context, id string, olds CommandState, news CommandArgs, preview bool) (CommandState, error) {
	state := CommandState{CommandArgs: news}
	if preview {
		return state, nil
	}
	var err error
	if news.Update != nil {
		state.Stdout, state.Stderr, _, err = state.run(ctx, *news.Update)
	} else if news.Create != nil {
		state.Stdout, state.Stderr, _, err = state.run(ctx, *news.Create)
	}
	return state, err
}

func (c *Command) Delete(ctx p.Context, id string, props CommandState) error {
	if props.Delete == nil {
		return nil
	}
	_, _, _, err := props.run(ctx, *props.Delete)
	return err
}

func (c *CommandState) run(ctx p.Context, command string) (string, string, string, error) {
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
	go util.CopyOutput(ctx, stdouttee, stdoutch, diag.Debug)
	go util.CopyOutput(ctx, stderrtee, stderrch, diag.Error)

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
