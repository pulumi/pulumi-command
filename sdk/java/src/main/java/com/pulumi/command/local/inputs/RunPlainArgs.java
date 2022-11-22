// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.command.local.inputs;

import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class RunPlainArgs extends com.pulumi.resources.InvokeArgs {

    public static final RunPlainArgs Empty = new RunPlainArgs();

    /**
     * A list of path globs to return as a single archive asset after the command completes.
     * 
     * When specifying glob patterns the following rules apply:
     * - We only include files not directories for assets and archives.
     * - Path separators are `/` on all platforms - including Windows.
     * - Patterns starting with `!` are &#39;exclude&#39; rules.
     * - Rules are evaluated in order, so exclude rules should be after inclusion rules.
     * - `*` matches anything except `/`
     * - `**` matches anything, _including_ `/`
     * - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
     * - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
     * 
     * #### Example
     * 
     * Given the rules:
     * 
     * When evaluating against this folder:
     * 
     * The following paths will be returned:
     * 
     */
    @Import(name="archivePaths")
    private @Nullable List<String> archivePaths;

    /**
     * @return A list of path globs to return as a single archive asset after the command completes.
     * 
     * When specifying glob patterns the following rules apply:
     * - We only include files not directories for assets and archives.
     * - Path separators are `/` on all platforms - including Windows.
     * - Patterns starting with `!` are &#39;exclude&#39; rules.
     * - Rules are evaluated in order, so exclude rules should be after inclusion rules.
     * - `*` matches anything except `/`
     * - `**` matches anything, _including_ `/`
     * - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
     * - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
     * 
     * #### Example
     * 
     * Given the rules:
     * 
     * When evaluating against this folder:
     * 
     * The following paths will be returned:
     * 
     */
    public Optional<List<String>> archivePaths() {
        return Optional.ofNullable(this.archivePaths);
    }

    /**
     * A list of path globs to read after the command completes.
     * 
     * When specifying glob patterns the following rules apply:
     * - We only include files not directories for assets and archives.
     * - Path separators are `/` on all platforms - including Windows.
     * - Patterns starting with `!` are &#39;exclude&#39; rules.
     * - Rules are evaluated in order, so exclude rules should be after inclusion rules.
     * - `*` matches anything except `/`
     * - `**` matches anything, _including_ `/`
     * - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
     * - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
     * 
     * #### Example
     * 
     * Given the rules:
     * 
     * When evaluating against this folder:
     * 
     * The following paths will be returned:
     * 
     */
    @Import(name="assetPaths")
    private @Nullable List<String> assetPaths;

    /**
     * @return A list of path globs to read after the command completes.
     * 
     * When specifying glob patterns the following rules apply:
     * - We only include files not directories for assets and archives.
     * - Path separators are `/` on all platforms - including Windows.
     * - Patterns starting with `!` are &#39;exclude&#39; rules.
     * - Rules are evaluated in order, so exclude rules should be after inclusion rules.
     * - `*` matches anything except `/`
     * - `**` matches anything, _including_ `/`
     * - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
     * - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
     * 
     * #### Example
     * 
     * Given the rules:
     * 
     * When evaluating against this folder:
     * 
     * The following paths will be returned:
     * 
     */
    public Optional<List<String>> assetPaths() {
        return Optional.ofNullable(this.assetPaths);
    }

    /**
     * The command to run.
     * 
     */
    @Import(name="command", required=true)
    private String command;

    /**
     * @return The command to run.
     * 
     */
    public String command() {
        return this.command;
    }

    /**
     * The directory from which to run the command from. If `dir` does not exist, then
     * `Command` will fail.
     * 
     */
    @Import(name="dir")
    private @Nullable String dir;

    /**
     * @return The directory from which to run the command from. If `dir` does not exist, then
     * `Command` will fail.
     * 
     */
    public Optional<String> dir() {
        return Optional.ofNullable(this.dir);
    }

    /**
     * Additional environment variables available to the command&#39;s process.
     * 
     */
    @Import(name="environment")
    private @Nullable Map<String,String> environment;

    /**
     * @return Additional environment variables available to the command&#39;s process.
     * 
     */
    public Optional<Map<String,String>> environment() {
        return Optional.ofNullable(this.environment);
    }

    /**
     * The program and arguments to run the command.
     * On Linux and macOS, defaults to: `[&#34;/bin/sh&#34;, &#34;-c&#34;]`. On Windows, defaults to: `[&#34;cmd&#34;, &#34;/C&#34;]`
     * 
     */
    @Import(name="interpreter")
    private @Nullable List<String> interpreter;

    /**
     * @return The program and arguments to run the command.
     * On Linux and macOS, defaults to: `[&#34;/bin/sh&#34;, &#34;-c&#34;]`. On Windows, defaults to: `[&#34;cmd&#34;, &#34;/C&#34;]`
     * 
     */
    public Optional<List<String>> interpreter() {
        return Optional.ofNullable(this.interpreter);
    }

    /**
     * Pass a string to the command&#39;s process as standard in
     * 
     */
    @Import(name="stdin")
    private @Nullable String stdin;

    /**
     * @return Pass a string to the command&#39;s process as standard in
     * 
     */
    public Optional<String> stdin() {
        return Optional.ofNullable(this.stdin);
    }

    private RunPlainArgs() {}

    private RunPlainArgs(RunPlainArgs $) {
        this.archivePaths = $.archivePaths;
        this.assetPaths = $.assetPaths;
        this.command = $.command;
        this.dir = $.dir;
        this.environment = $.environment;
        this.interpreter = $.interpreter;
        this.stdin = $.stdin;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(RunPlainArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private RunPlainArgs $;

        public Builder() {
            $ = new RunPlainArgs();
        }

        public Builder(RunPlainArgs defaults) {
            $ = new RunPlainArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param archivePaths A list of path globs to return as a single archive asset after the command completes.
         * 
         * When specifying glob patterns the following rules apply:
         * - We only include files not directories for assets and archives.
         * - Path separators are `/` on all platforms - including Windows.
         * - Patterns starting with `!` are &#39;exclude&#39; rules.
         * - Rules are evaluated in order, so exclude rules should be after inclusion rules.
         * - `*` matches anything except `/`
         * - `**` matches anything, _including_ `/`
         * - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
         * - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
         * 
         * #### Example
         * 
         * Given the rules:
         * 
         * When evaluating against this folder:
         * 
         * The following paths will be returned:
         * 
         * @return builder
         * 
         */
        public Builder archivePaths(@Nullable List<String> archivePaths) {
            $.archivePaths = archivePaths;
            return this;
        }

        /**
         * @param archivePaths A list of path globs to return as a single archive asset after the command completes.
         * 
         * When specifying glob patterns the following rules apply:
         * - We only include files not directories for assets and archives.
         * - Path separators are `/` on all platforms - including Windows.
         * - Patterns starting with `!` are &#39;exclude&#39; rules.
         * - Rules are evaluated in order, so exclude rules should be after inclusion rules.
         * - `*` matches anything except `/`
         * - `**` matches anything, _including_ `/`
         * - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
         * - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
         * 
         * #### Example
         * 
         * Given the rules:
         * 
         * When evaluating against this folder:
         * 
         * The following paths will be returned:
         * 
         * @return builder
         * 
         */
        public Builder archivePaths(String... archivePaths) {
            return archivePaths(List.of(archivePaths));
        }

        /**
         * @param assetPaths A list of path globs to read after the command completes.
         * 
         * When specifying glob patterns the following rules apply:
         * - We only include files not directories for assets and archives.
         * - Path separators are `/` on all platforms - including Windows.
         * - Patterns starting with `!` are &#39;exclude&#39; rules.
         * - Rules are evaluated in order, so exclude rules should be after inclusion rules.
         * - `*` matches anything except `/`
         * - `**` matches anything, _including_ `/`
         * - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
         * - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
         * 
         * #### Example
         * 
         * Given the rules:
         * 
         * When evaluating against this folder:
         * 
         * The following paths will be returned:
         * 
         * @return builder
         * 
         */
        public Builder assetPaths(@Nullable List<String> assetPaths) {
            $.assetPaths = assetPaths;
            return this;
        }

        /**
         * @param assetPaths A list of path globs to read after the command completes.
         * 
         * When specifying glob patterns the following rules apply:
         * - We only include files not directories for assets and archives.
         * - Path separators are `/` on all platforms - including Windows.
         * - Patterns starting with `!` are &#39;exclude&#39; rules.
         * - Rules are evaluated in order, so exclude rules should be after inclusion rules.
         * - `*` matches anything except `/`
         * - `**` matches anything, _including_ `/`
         * - All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
         * - For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)
         * 
         * #### Example
         * 
         * Given the rules:
         * 
         * When evaluating against this folder:
         * 
         * The following paths will be returned:
         * 
         * @return builder
         * 
         */
        public Builder assetPaths(String... assetPaths) {
            return assetPaths(List.of(assetPaths));
        }

        /**
         * @param command The command to run.
         * 
         * @return builder
         * 
         */
        public Builder command(String command) {
            $.command = command;
            return this;
        }

        /**
         * @param dir The directory from which to run the command from. If `dir` does not exist, then
         * `Command` will fail.
         * 
         * @return builder
         * 
         */
        public Builder dir(@Nullable String dir) {
            $.dir = dir;
            return this;
        }

        /**
         * @param environment Additional environment variables available to the command&#39;s process.
         * 
         * @return builder
         * 
         */
        public Builder environment(@Nullable Map<String,String> environment) {
            $.environment = environment;
            return this;
        }

        /**
         * @param interpreter The program and arguments to run the command.
         * On Linux and macOS, defaults to: `[&#34;/bin/sh&#34;, &#34;-c&#34;]`. On Windows, defaults to: `[&#34;cmd&#34;, &#34;/C&#34;]`
         * 
         * @return builder
         * 
         */
        public Builder interpreter(@Nullable List<String> interpreter) {
            $.interpreter = interpreter;
            return this;
        }

        /**
         * @param interpreter The program and arguments to run the command.
         * On Linux and macOS, defaults to: `[&#34;/bin/sh&#34;, &#34;-c&#34;]`. On Windows, defaults to: `[&#34;cmd&#34;, &#34;/C&#34;]`
         * 
         * @return builder
         * 
         */
        public Builder interpreter(String... interpreter) {
            return interpreter(List.of(interpreter));
        }

        /**
         * @param stdin Pass a string to the command&#39;s process as standard in
         * 
         * @return builder
         * 
         */
        public Builder stdin(@Nullable String stdin) {
            $.stdin = stdin;
            return this;
        }

        public RunPlainArgs build() {
            $.command = Objects.requireNonNull($.command, "expected parameter 'command' to be non-null");
            return $;
        }
    }

}
