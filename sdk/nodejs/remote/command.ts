// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../types/input";
import * as outputs from "../types/output";
import * as enums from "../types/enums";
import * as utilities from "../utilities";

/**
 * A command to run on a remote host. The connection is established via ssh.
 *
 * ## Example Usage
 *
 * ### A Basic Example
 * This program connects to a server and runs the `hostname` command. The output is then available via the `stdout` property.
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as command from "@pulumi/command";
 *
 * const config = new pulumi.Config();
 * const server = config.require("server");
 * const userName = config.require("userName");
 * const privateKey = config.require("privateKey");
 *
 * const hostnameCmd = new command.remote.Command("hostnameCmd", {
 *     create: "hostname",
 *     connection: {
 *         host: server,
 *         user: userName,
 *         privateKey: privateKey,
 *     },
 * });
 * export const hostname = hostnameCmd.stdout;
 * ```
 *
 * ### Triggers
 * This example defines several trigger values of various kinds. Changes to any of them will cause `cmd` to be re-run.
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as command from "@pulumi/command";
 * import * as random from "@pulumi/random";
 *
 * const str = "foo";
 * const fileAsset = new pulumi.asset.FileAsset("Pulumi.yaml");
 * const rand = new random.RandomString("rand", {length: 5});
 * const localFile = new command.local.Command("localFile", {
 *     create: "touch foo.txt",
 *     archivePaths: ["*.txt"],
 * });
 * const cmd = new command.remote.Command("cmd", {
 *     connection: {
 *         host: "insert host here",
 *     },
 *     create: "echo create > op.txt",
 *     delete: "echo delete >> op.txt",
 *     triggers: [
 *         str,
 *         rand.result,
 *         fileAsset,
 *         localFile.archive,
 *     ],
 * });
 *
 * ```
 */
export class Command extends pulumi.CustomResource {
    /**
     * Get an existing Command resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    public static get(name: string, id: pulumi.Input<pulumi.ID>, opts?: pulumi.CustomResourceOptions): Command {
        return new Command(name, undefined as any, { ...opts, id: id });
    }

    /** @internal */
    public static readonly __pulumiType = 'command:remote:Command';

    /**
     * Returns true if the given object is an instance of Command.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is Command {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === Command.__pulumiType;
    }

    /**
     * If the previous command's stdout and stderr (as generated by the prior create/update) is
     * injected into the environment of the next run as PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR.
     * Defaults to true.
     */
    public readonly addPreviousOutputInEnv!: pulumi.Output<boolean | undefined>;
    /**
     * The parameters with which to connect to the remote host.
     */
    public readonly connection!: pulumi.Output<outputs.remote.Connection>;
    /**
     * The command to run once on resource creation.
     *
     * If an `update` command isn't provided, then `create` will also be run when the resource's inputs are modified.
     *
     * Note that this command will not be executed if the resource has already been created and its inputs are unchanged.
     *
     * Use `local.runOutput` if you need to run a command on every execution of your program.
     */
    public readonly create!: pulumi.Output<string | undefined>;
    /**
     * The command to run on resource delettion.
     *
     * The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the stdout and stderr properties of the Command resource from previous create or update steps.
     */
    public readonly delete!: pulumi.Output<string | undefined>;
    /**
     * Additional environment variables available to the command's process.
     * Note that this only works if the SSH server is configured to accept these variables via AcceptEnv.
     * Alternatively, if a Bash-like shell runs the command on the remote host, you could prefix the command itself
     * with the variables in the form 'VAR=value command'.
     */
    public readonly environment!: pulumi.Output<{[key: string]: string} | undefined>;
    /**
     * If the command's stdout and stderr should be logged. This doesn't affect the capturing of
     * stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
     * outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.
     */
    public readonly logging!: pulumi.Output<enums.remote.Logging | undefined>;
    /**
     * The standard error of the command's process
     */
    public /*out*/ readonly stderr!: pulumi.Output<string>;
    /**
     * Pass a string to the command's process as standard in
     */
    public readonly stdin!: pulumi.Output<string | undefined>;
    /**
     * The standard output of the command's process
     */
    public /*out*/ readonly stdout!: pulumi.Output<string>;
    /**
     * The resource will be updated (or replaced) if any of these values change.
     *
     * The trigger values can be of any type.
     *
     * If the `update` command was provided the resource will be updated, otherwise it will be replaced using the `create` command.
     *
     * Please see the resource documentation for examples.
     */
    public readonly triggers!: pulumi.Output<any[] | undefined>;
    /**
     * The command to run when the resource is updated.
     *
     * If empty, the create command will be executed instead.
     *
     * Note that this command will not run if the resource's inputs are unchanged.
     *
     * Use `local.runOutput` if you need to run a command on every execution of your program.
     *
     * The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the `stdout` and `stderr` properties of the Command resource from previous create or update steps.
     */
    public readonly update!: pulumi.Output<string | undefined>;

    /**
     * Create a Command resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: CommandArgs, opts?: pulumi.CustomResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            if ((!args || args.connection === undefined) && !opts.urn) {
                throw new Error("Missing required property 'connection'");
            }
            resourceInputs["addPreviousOutputInEnv"] = (args ? args.addPreviousOutputInEnv : undefined) ?? true;
            resourceInputs["connection"] = args?.connection ? pulumi.secret((args.connection ? pulumi.output(args.connection).apply(inputs.remote.connectionArgsProvideDefaults) : undefined)) : undefined;
            resourceInputs["create"] = args ? args.create : undefined;
            resourceInputs["delete"] = args ? args.delete : undefined;
            resourceInputs["environment"] = args ? args.environment : undefined;
            resourceInputs["logging"] = args ? args.logging : undefined;
            resourceInputs["stdin"] = args ? args.stdin : undefined;
            resourceInputs["triggers"] = args ? args.triggers : undefined;
            resourceInputs["update"] = args ? args.update : undefined;
            resourceInputs["stderr"] = undefined /*out*/;
            resourceInputs["stdout"] = undefined /*out*/;
        } else {
            resourceInputs["addPreviousOutputInEnv"] = undefined /*out*/;
            resourceInputs["connection"] = undefined /*out*/;
            resourceInputs["create"] = undefined /*out*/;
            resourceInputs["delete"] = undefined /*out*/;
            resourceInputs["environment"] = undefined /*out*/;
            resourceInputs["logging"] = undefined /*out*/;
            resourceInputs["stderr"] = undefined /*out*/;
            resourceInputs["stdin"] = undefined /*out*/;
            resourceInputs["stdout"] = undefined /*out*/;
            resourceInputs["triggers"] = undefined /*out*/;
            resourceInputs["update"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        const secretOpts = { additionalSecretOutputs: ["connection"] };
        opts = pulumi.mergeOptions(opts, secretOpts);
        const replaceOnChanges = { replaceOnChanges: ["triggers[*]"] };
        opts = pulumi.mergeOptions(opts, replaceOnChanges);
        super(Command.__pulumiType, name, resourceInputs, opts);
    }
}

/**
 * The set of arguments for constructing a Command resource.
 */
export interface CommandArgs {
    /**
     * If the previous command's stdout and stderr (as generated by the prior create/update) is
     * injected into the environment of the next run as PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR.
     * Defaults to true.
     */
    addPreviousOutputInEnv?: pulumi.Input<boolean>;
    /**
     * The parameters with which to connect to the remote host.
     */
    connection: pulumi.Input<inputs.remote.ConnectionArgs>;
    /**
     * The command to run once on resource creation.
     *
     * If an `update` command isn't provided, then `create` will also be run when the resource's inputs are modified.
     *
     * Note that this command will not be executed if the resource has already been created and its inputs are unchanged.
     *
     * Use `local.runOutput` if you need to run a command on every execution of your program.
     */
    create?: pulumi.Input<string>;
    /**
     * The command to run on resource delettion.
     *
     * The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the stdout and stderr properties of the Command resource from previous create or update steps.
     */
    delete?: pulumi.Input<string>;
    /**
     * Additional environment variables available to the command's process.
     * Note that this only works if the SSH server is configured to accept these variables via AcceptEnv.
     * Alternatively, if a Bash-like shell runs the command on the remote host, you could prefix the command itself
     * with the variables in the form 'VAR=value command'.
     */
    environment?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * If the command's stdout and stderr should be logged. This doesn't affect the capturing of
     * stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
     * outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.
     */
    logging?: pulumi.Input<enums.remote.Logging>;
    /**
     * Pass a string to the command's process as standard in
     */
    stdin?: pulumi.Input<string>;
    /**
     * The resource will be updated (or replaced) if any of these values change.
     *
     * The trigger values can be of any type.
     *
     * If the `update` command was provided the resource will be updated, otherwise it will be replaced using the `create` command.
     *
     * Please see the resource documentation for examples.
     */
    triggers?: pulumi.Input<any[]>;
    /**
     * The command to run when the resource is updated.
     *
     * If empty, the create command will be executed instead.
     *
     * Note that this command will not run if the resource's inputs are unchanged.
     *
     * Use `local.runOutput` if you need to run a command on every execution of your program.
     *
     * The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the `stdout` and `stderr` properties of the Command resource from previous create or update steps.
     */
    update?: pulumi.Input<string>;
}
