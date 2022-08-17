// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.command.local.outputs;

import com.pulumi.asset.Archive;
import com.pulumi.asset.AssetOrArchive;
import com.pulumi.core.annotations.CustomType;
import java.lang.String;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class RunResult {
    /**
     * @return An archive asset containing files found after running the command.
     * 
     */
    private final @Nullable Archive archive;
    /**
     * @return A map of assets found after running the command.
     * The key is the relative path from the command dir
     * 
     */
    private final @Nullable Map<String,AssetOrArchive> assets;
    /**
     * @return The command to run.
     * 
     */
    private final String command;
    /**
     * @return The directory from which the command was run from.
     * 
     */
    private final @Nullable String dir;
    /**
     * @return Additional environment variables available to the command&#39;s process.
     * 
     */
    private final @Nullable Map<String,String> environment;
    /**
     * @return The program and arguments to run the command.
     * For example: `[&#34;/bin/sh&#34;, &#34;-c&#34;]`
     * 
     */
    private final @Nullable List<String> interpreter;
    /**
     * @return The standard error of the command&#39;s process
     * 
     */
    private final String stderr;
    /**
     * @return String passed to the command&#39;s process as standard in.
     * 
     */
    private final String stdin;
    /**
     * @return The standard output of the command&#39;s process
     * 
     */
    private final @Nullable String stdout;

    @CustomType.Constructor
    private RunResult(
        @CustomType.Parameter("archive") @Nullable Archive archive,
        @CustomType.Parameter("assets") @Nullable Map<String,AssetOrArchive> assets,
        @CustomType.Parameter("command") String command,
        @CustomType.Parameter("dir") @Nullable String dir,
        @CustomType.Parameter("environment") @Nullable Map<String,String> environment,
        @CustomType.Parameter("interpreter") @Nullable List<String> interpreter,
        @CustomType.Parameter("stderr") String stderr,
        @CustomType.Parameter("stdin") String stdin,
        @CustomType.Parameter("stdout") @Nullable String stdout) {
        this.archive = archive;
        this.assets = assets;
        this.command = command;
        this.dir = dir;
        this.environment = environment;
        this.interpreter = interpreter;
        this.stderr = stderr;
        this.stdin = stdin;
        this.stdout = stdout;
    }

    /**
     * @return An archive asset containing files found after running the command.
     * 
     */
    public Optional<Archive> archive() {
        return Optional.ofNullable(this.archive);
    }
    /**
     * @return A map of assets found after running the command.
     * The key is the relative path from the command dir
     * 
     */
    public Map<String,AssetOrArchive> assets() {
        return this.assets == null ? Map.of() : this.assets;
    }
    /**
     * @return The command to run.
     * 
     */
    public String command() {
        return this.command;
    }
    /**
     * @return The directory from which the command was run from.
     * 
     */
    public Optional<String> dir() {
        return Optional.ofNullable(this.dir);
    }
    /**
     * @return Additional environment variables available to the command&#39;s process.
     * 
     */
    public Map<String,String> environment() {
        return this.environment == null ? Map.of() : this.environment;
    }
    /**
     * @return The program and arguments to run the command.
     * For example: `[&#34;/bin/sh&#34;, &#34;-c&#34;]`
     * 
     */
    public List<String> interpreter() {
        return this.interpreter == null ? List.of() : this.interpreter;
    }
    /**
     * @return The standard error of the command&#39;s process
     * 
     */
    public String stderr() {
        return this.stderr;
    }
    /**
     * @return String passed to the command&#39;s process as standard in.
     * 
     */
    public String stdin() {
        return this.stdin;
    }
    /**
     * @return The standard output of the command&#39;s process
     * 
     */
    public Optional<String> stdout() {
        return Optional.ofNullable(this.stdout);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(RunResult defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private @Nullable Archive archive;
        private @Nullable Map<String,AssetOrArchive> assets;
        private String command;
        private @Nullable String dir;
        private @Nullable Map<String,String> environment;
        private @Nullable List<String> interpreter;
        private String stderr;
        private String stdin;
        private @Nullable String stdout;

        public Builder() {
    	      // Empty
        }

        public Builder(RunResult defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.archive = defaults.archive;
    	      this.assets = defaults.assets;
    	      this.command = defaults.command;
    	      this.dir = defaults.dir;
    	      this.environment = defaults.environment;
    	      this.interpreter = defaults.interpreter;
    	      this.stderr = defaults.stderr;
    	      this.stdin = defaults.stdin;
    	      this.stdout = defaults.stdout;
        }

        public Builder archive(@Nullable Archive archive) {
            this.archive = archive;
            return this;
        }
        public Builder assets(@Nullable Map<String,AssetOrArchive> assets) {
            this.assets = assets;
            return this;
        }
        public Builder command(String command) {
            this.command = Objects.requireNonNull(command);
            return this;
        }
        public Builder dir(@Nullable String dir) {
            this.dir = dir;
            return this;
        }
        public Builder environment(@Nullable Map<String,String> environment) {
            this.environment = environment;
            return this;
        }
        public Builder interpreter(@Nullable List<String> interpreter) {
            this.interpreter = interpreter;
            return this;
        }
        public Builder interpreter(String... interpreter) {
            return interpreter(List.of(interpreter));
        }
        public Builder stderr(String stderr) {
            this.stderr = Objects.requireNonNull(stderr);
            return this;
        }
        public Builder stdin(String stdin) {
            this.stdin = Objects.requireNonNull(stdin);
            return this;
        }
        public Builder stdout(@Nullable String stdout) {
            this.stdout = stdout;
            return this;
        }        public RunResult build() {
            return new RunResult(archive, assets, command, dir, environment, interpreter, stderr, stdin, stdout);
        }
    }
}
