// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Command.Local
{
    public static class Run
    {
        /// <summary>
        /// A local command to be executed.
        /// This command will always be run on any preview or deployment. Use `local.Command` to avoid duplicating executions.
        /// </summary>
        public static Task<RunResult> InvokeAsync(RunArgs args, InvokeOptions? options = null)
            => global::Pulumi.Deployment.Instance.InvokeAsync<RunResult>("command:local:run", args ?? new RunArgs(), options.WithDefaults());

        /// <summary>
        /// A local command to be executed.
        /// This command will always be run on any preview or deployment. Use `local.Command` to avoid duplicating executions.
        /// </summary>
        public static Output<RunResult> Invoke(RunInvokeArgs args, InvokeOptions? options = null)
            => global::Pulumi.Deployment.Instance.Invoke<RunResult>("command:local:run", args ?? new RunInvokeArgs(), options.WithDefaults());
    }


    public sealed class RunArgs : global::Pulumi.InvokeArgs
    {
        /// <summary>
        /// If the previous command's stdout and stderr (as generated by the prior create/update) is
        /// injected into the environment of the next run as PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR.
        /// Defaults to true.
        /// </summary>
        [Input("addPreviousOutputInEnv")]
        public bool? AddPreviousOutputInEnv { get; set; }

        [Input("archivePaths")]
        private List<string>? _archivePaths;

        /// <summary>
        /// A list of path globs to return as a single archive asset after the command completes.
        /// 
        /// When specifying glob patterns the following rules apply:
        /// - We only include files not directories for assets and archives.
        /// - Path separators are `/` on all platforms - including Windows.
        /// - Patterns starting with `!` are 'exclude' rules.
        /// - Rules are evaluated in order, so exclude rules should be after inclusion rules.
        /// - `*` matches anything except `/`
        /// - `**` matches anything, _including_ `/`
        /// - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
        /// - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
        /// 
        /// #### Example
        /// 
        /// Given the rules:
        /// ```yaml
        /// - "assets/**"
        /// - "src/**.js"
        /// - "!**secret.*"
        /// ```
        /// 
        /// When evaluating against this folder:
        /// 
        /// ```yaml
        /// - assets/
        ///   - logos/
        ///     - logo.svg
        /// - src/
        ///   - index.js
        ///   - secret.js
        /// ```
        /// 
        /// The following paths will be returned:
        /// 
        /// ```yaml
        /// - assets/logos/logo.svg
        /// - src/index.js
        /// ```
        /// </summary>
        public List<string> ArchivePaths
        {
            get => _archivePaths ?? (_archivePaths = new List<string>());
            set => _archivePaths = value;
        }

        [Input("assetPaths")]
        private List<string>? _assetPaths;

        /// <summary>
        /// A list of path globs to read after the command completes.
        /// 
        /// When specifying glob patterns the following rules apply:
        /// - We only include files not directories for assets and archives.
        /// - Path separators are `/` on all platforms - including Windows.
        /// - Patterns starting with `!` are 'exclude' rules.
        /// - Rules are evaluated in order, so exclude rules should be after inclusion rules.
        /// - `*` matches anything except `/`
        /// - `**` matches anything, _including_ `/`
        /// - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
        /// - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
        /// 
        /// #### Example
        /// 
        /// Given the rules:
        /// ```yaml
        /// - "assets/**"
        /// - "src/**.js"
        /// - "!**secret.*"
        /// ```
        /// 
        /// When evaluating against this folder:
        /// 
        /// ```yaml
        /// - assets/
        ///   - logos/
        ///     - logo.svg
        /// - src/
        ///   - index.js
        ///   - secret.js
        /// ```
        /// 
        /// The following paths will be returned:
        /// 
        /// ```yaml
        /// - assets/logos/logo.svg
        /// - src/index.js
        /// ```
        /// </summary>
        public List<string> AssetPaths
        {
            get => _assetPaths ?? (_assetPaths = new List<string>());
            set => _assetPaths = value;
        }

        /// <summary>
        /// The command to run.
        /// </summary>
        [Input("command", required: true)]
        public string Command { get; set; } = null!;

        /// <summary>
        /// The directory from which to run the command from. If `dir` does not exist, then
        /// `Command` will fail.
        /// </summary>
        [Input("dir")]
        public string? Dir { get; set; }

        [Input("environment")]
        private Dictionary<string, string>? _environment;

        /// <summary>
        /// Additional environment variables available to the command's process.
        /// </summary>
        public Dictionary<string, string> Environment
        {
            get => _environment ?? (_environment = new Dictionary<string, string>());
            set => _environment = value;
        }

        [Input("interpreter")]
        private List<string>? _interpreter;

        /// <summary>
        /// The program and arguments to run the command.
        /// On Linux and macOS, defaults to: `["/bin/sh", "-c"]`. On Windows, defaults to: `["cmd", "/C"]`
        /// </summary>
        public List<string> Interpreter
        {
            get => _interpreter ?? (_interpreter = new List<string>());
            set => _interpreter = value;
        }

        /// <summary>
        /// If the command's stdout and stderr should be logged. This doesn't affect the capturing of
        /// stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
        /// outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.
        /// </summary>
        [Input("logging")]
        public Pulumi.Command.Common.Logging? Logging { get; set; }

        /// <summary>
        /// Pass a string to the command's process as standard in
        /// </summary>
        [Input("stdin")]
        public string? Stdin { get; set; }

        public RunArgs()
        {
            AddPreviousOutputInEnv = true;
        }
        public static new RunArgs Empty => new RunArgs();
    }

    public sealed class RunInvokeArgs : global::Pulumi.InvokeArgs
    {
        /// <summary>
        /// If the previous command's stdout and stderr (as generated by the prior create/update) is
        /// injected into the environment of the next run as PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR.
        /// Defaults to true.
        /// </summary>
        [Input("addPreviousOutputInEnv")]
        public Input<bool>? AddPreviousOutputInEnv { get; set; }

        [Input("archivePaths")]
        private InputList<string>? _archivePaths;

        /// <summary>
        /// A list of path globs to return as a single archive asset after the command completes.
        /// 
        /// When specifying glob patterns the following rules apply:
        /// - We only include files not directories for assets and archives.
        /// - Path separators are `/` on all platforms - including Windows.
        /// - Patterns starting with `!` are 'exclude' rules.
        /// - Rules are evaluated in order, so exclude rules should be after inclusion rules.
        /// - `*` matches anything except `/`
        /// - `**` matches anything, _including_ `/`
        /// - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
        /// - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
        /// 
        /// #### Example
        /// 
        /// Given the rules:
        /// ```yaml
        /// - "assets/**"
        /// - "src/**.js"
        /// - "!**secret.*"
        /// ```
        /// 
        /// When evaluating against this folder:
        /// 
        /// ```yaml
        /// - assets/
        ///   - logos/
        ///     - logo.svg
        /// - src/
        ///   - index.js
        ///   - secret.js
        /// ```
        /// 
        /// The following paths will be returned:
        /// 
        /// ```yaml
        /// - assets/logos/logo.svg
        /// - src/index.js
        /// ```
        /// </summary>
        public InputList<string> ArchivePaths
        {
            get => _archivePaths ?? (_archivePaths = new InputList<string>());
            set => _archivePaths = value;
        }

        [Input("assetPaths")]
        private InputList<string>? _assetPaths;

        /// <summary>
        /// A list of path globs to read after the command completes.
        /// 
        /// When specifying glob patterns the following rules apply:
        /// - We only include files not directories for assets and archives.
        /// - Path separators are `/` on all platforms - including Windows.
        /// - Patterns starting with `!` are 'exclude' rules.
        /// - Rules are evaluated in order, so exclude rules should be after inclusion rules.
        /// - `*` matches anything except `/`
        /// - `**` matches anything, _including_ `/`
        /// - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
        /// - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
        /// 
        /// #### Example
        /// 
        /// Given the rules:
        /// ```yaml
        /// - "assets/**"
        /// - "src/**.js"
        /// - "!**secret.*"
        /// ```
        /// 
        /// When evaluating against this folder:
        /// 
        /// ```yaml
        /// - assets/
        ///   - logos/
        ///     - logo.svg
        /// - src/
        ///   - index.js
        ///   - secret.js
        /// ```
        /// 
        /// The following paths will be returned:
        /// 
        /// ```yaml
        /// - assets/logos/logo.svg
        /// - src/index.js
        /// ```
        /// </summary>
        public InputList<string> AssetPaths
        {
            get => _assetPaths ?? (_assetPaths = new InputList<string>());
            set => _assetPaths = value;
        }

        /// <summary>
        /// The command to run.
        /// </summary>
        [Input("command", required: true)]
        public Input<string> Command { get; set; } = null!;

        /// <summary>
        /// The directory from which to run the command from. If `dir` does not exist, then
        /// `Command` will fail.
        /// </summary>
        [Input("dir")]
        public Input<string>? Dir { get; set; }

        [Input("environment")]
        private InputMap<string>? _environment;

        /// <summary>
        /// Additional environment variables available to the command's process.
        /// </summary>
        public InputMap<string> Environment
        {
            get => _environment ?? (_environment = new InputMap<string>());
            set => _environment = value;
        }

        [Input("interpreter")]
        private InputList<string>? _interpreter;

        /// <summary>
        /// The program and arguments to run the command.
        /// On Linux and macOS, defaults to: `["/bin/sh", "-c"]`. On Windows, defaults to: `["cmd", "/C"]`
        /// </summary>
        public InputList<string> Interpreter
        {
            get => _interpreter ?? (_interpreter = new InputList<string>());
            set => _interpreter = value;
        }

        /// <summary>
        /// If the command's stdout and stderr should be logged. This doesn't affect the capturing of
        /// stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
        /// outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.
        /// </summary>
        [Input("logging")]
        public Input<Pulumi.Command.Common.Logging>? Logging { get; set; }

        /// <summary>
        /// Pass a string to the command's process as standard in
        /// </summary>
        [Input("stdin")]
        public Input<string>? Stdin { get; set; }

        public RunInvokeArgs()
        {
            AddPreviousOutputInEnv = true;
        }
        public static new RunInvokeArgs Empty => new RunInvokeArgs();
    }


    [OutputType]
    public sealed class RunResult
    {
        /// <summary>
        /// If the previous command's stdout and stderr (as generated by the prior create/update) is
        /// injected into the environment of the next run as PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR.
        /// Defaults to true.
        /// </summary>
        public readonly bool? AddPreviousOutputInEnv;
        /// <summary>
        /// An archive asset containing files found after running the command.
        /// </summary>
        public readonly Archive? Archive;
        /// <summary>
        /// A list of path globs to return as a single archive asset after the command completes.
        /// 
        /// When specifying glob patterns the following rules apply:
        /// - We only include files not directories for assets and archives.
        /// - Path separators are `/` on all platforms - including Windows.
        /// - Patterns starting with `!` are 'exclude' rules.
        /// - Rules are evaluated in order, so exclude rules should be after inclusion rules.
        /// - `*` matches anything except `/`
        /// - `**` matches anything, _including_ `/`
        /// - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
        /// - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
        /// 
        /// #### Example
        /// 
        /// Given the rules:
        /// ```yaml
        /// - "assets/**"
        /// - "src/**.js"
        /// - "!**secret.*"
        /// ```
        /// 
        /// When evaluating against this folder:
        /// 
        /// ```yaml
        /// - assets/
        ///   - logos/
        ///     - logo.svg
        /// - src/
        ///   - index.js
        ///   - secret.js
        /// ```
        /// 
        /// The following paths will be returned:
        /// 
        /// ```yaml
        /// - assets/logos/logo.svg
        /// - src/index.js
        /// ```
        /// </summary>
        public readonly ImmutableArray<string> ArchivePaths;
        /// <summary>
        /// A list of path globs to read after the command completes.
        /// 
        /// When specifying glob patterns the following rules apply:
        /// - We only include files not directories for assets and archives.
        /// - Path separators are `/` on all platforms - including Windows.
        /// - Patterns starting with `!` are 'exclude' rules.
        /// - Rules are evaluated in order, so exclude rules should be after inclusion rules.
        /// - `*` matches anything except `/`
        /// - `**` matches anything, _including_ `/`
        /// - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
        /// - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
        /// 
        /// #### Example
        /// 
        /// Given the rules:
        /// ```yaml
        /// - "assets/**"
        /// - "src/**.js"
        /// - "!**secret.*"
        /// ```
        /// 
        /// When evaluating against this folder:
        /// 
        /// ```yaml
        /// - assets/
        ///   - logos/
        ///     - logo.svg
        /// - src/
        ///   - index.js
        ///   - secret.js
        /// ```
        /// 
        /// The following paths will be returned:
        /// 
        /// ```yaml
        /// - assets/logos/logo.svg
        /// - src/index.js
        /// ```
        /// </summary>
        public readonly ImmutableArray<string> AssetPaths;
        /// <summary>
        /// A map of assets found after running the command.
        /// The key is the relative path from the command dir
        /// </summary>
        public readonly ImmutableDictionary<string, AssetOrArchive>? Assets;
        /// <summary>
        /// The command to run.
        /// </summary>
        public readonly string Command;
        /// <summary>
        /// The directory from which to run the command from. If `dir` does not exist, then
        /// `Command` will fail.
        /// </summary>
        public readonly string? Dir;
        /// <summary>
        /// Additional environment variables available to the command's process.
        /// </summary>
        public readonly ImmutableDictionary<string, string>? Environment;
        /// <summary>
        /// The program and arguments to run the command.
        /// On Linux and macOS, defaults to: `["/bin/sh", "-c"]`. On Windows, defaults to: `["cmd", "/C"]`
        /// </summary>
        public readonly ImmutableArray<string> Interpreter;
        /// <summary>
        /// If the command's stdout and stderr should be logged. This doesn't affect the capturing of
        /// stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
        /// outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.
        /// </summary>
        public readonly Pulumi.Command.Common.Logging? Logging;
        /// <summary>
        /// The standard error of the command's process
        /// </summary>
        public readonly string Stderr;
        /// <summary>
        /// Pass a string to the command's process as standard in
        /// </summary>
        public readonly string? Stdin;
        /// <summary>
        /// The standard output of the command's process
        /// </summary>
        public readonly string Stdout;

        [OutputConstructor]
        private RunResult(
            bool? addPreviousOutputInEnv,

            Archive? archive,

            ImmutableArray<string> archivePaths,

            ImmutableArray<string> assetPaths,

            ImmutableDictionary<string, AssetOrArchive>? assets,

            string command,

            string? dir,

            ImmutableDictionary<string, string>? environment,

            ImmutableArray<string> interpreter,

            Pulumi.Command.Common.Logging? logging,

            string stderr,

            string? stdin,

            string stdout)
        {
            AddPreviousOutputInEnv = addPreviousOutputInEnv;
            Archive = archive;
            ArchivePaths = archivePaths;
            AssetPaths = assetPaths;
            Assets = assets;
            Command = command;
            Dir = dir;
            Environment = environment;
            Interpreter = interpreter;
            Logging = logging;
            Stderr = stderr;
            Stdin = stdin;
            Stdout = stdout;
        }
    }
}
