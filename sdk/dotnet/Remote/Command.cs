// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Command.Remote
{
    /// <summary>
    /// A command to run on a remote host.
    /// The connection is established via ssh.
    /// </summary>
    [CommandResourceType("command:remote:Command")]
    public partial class Command : global::Pulumi.CustomResource
    {
        /// <summary>
        /// The parameters with which to connect to the remote host.
        /// </summary>
        [Output("connection")]
        public Output<Outputs.Connection> Connection { get; private set; } = null!;

        /// <summary>
        /// The command to run on create.
        /// </summary>
        [Output("create")]
        public Output<string?> Create { get; private set; } = null!;

        /// <summary>
        /// The command to run on delete.
        /// </summary>
        [Output("delete")]
        public Output<string?> Delete { get; private set; } = null!;

        /// <summary>
        /// Additional environment variables available to the command's process.
        /// </summary>
        [Output("environment")]
        public Output<ImmutableDictionary<string, string>?> Environment { get; private set; } = null!;

        /// <summary>
        /// The standard error of the command's process
        /// </summary>
        [Output("stderr")]
        public Output<string> Stderr { get; private set; } = null!;

        /// <summary>
        /// Pass a string to the command's process as standard in
        /// </summary>
        [Output("stdin")]
        public Output<string?> Stdin { get; private set; } = null!;

        /// <summary>
        /// The standard output of the command's process
        /// </summary>
        [Output("stdout")]
        public Output<string> Stdout { get; private set; } = null!;

        /// <summary>
        /// Trigger replacements on changes to this input.
        /// </summary>
        [Output("triggers")]
        public Output<ImmutableArray<object>> Triggers { get; private set; } = null!;

        /// <summary>
        /// The command to run on update, if empty, create will run again.
        /// </summary>
        [Output("update")]
        public Output<string?> Update { get; private set; } = null!;


        /// <summary>
        /// Create a Command resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Command(string name, CommandArgs args, CustomResourceOptions? options = null)
            : base("command:remote:Command", name, args ?? new CommandArgs(), MakeResourceOptions(options, ""))
        {
        }

        private Command(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("command:remote:Command", name, null, MakeResourceOptions(options, id))
        {
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
                AdditionalSecretOutputs =
                {
                    "connection",
                },
                ReplaceOnChanges =
                {
                    "triggers[*]",
                },
            };
            var merged = CustomResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
        /// <summary>
        /// Get an existing Command resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static Command Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new Command(name, id, options);
        }
    }

    public sealed class CommandArgs : global::Pulumi.ResourceArgs
    {
        [Input("connection", required: true)]
        private Input<Inputs.ConnectionArgs>? _connection;

        /// <summary>
        /// The parameters with which to connect to the remote host.
        /// </summary>
        public Input<Inputs.ConnectionArgs>? Connection
        {
            get => _connection;
            set
            {
                var emptySecret = Output.CreateSecret(0);
                _connection = Output.Tuple<Input<Inputs.ConnectionArgs>?, int>(value, emptySecret).Apply(t => t.Item1);
            }
        }

        /// <summary>
        /// The command to run on create.
        /// </summary>
        [Input("create")]
        public Input<string>? Create { get; set; }

        /// <summary>
        /// The command to run on delete.
        /// </summary>
        [Input("delete")]
        public Input<string>? Delete { get; set; }

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

        /// <summary>
        /// Pass a string to the command's process as standard in
        /// </summary>
        [Input("stdin")]
        public Input<string>? Stdin { get; set; }

        [Input("triggers")]
        private InputList<object>? _triggers;

        /// <summary>
        /// Trigger replacements on changes to this input.
        /// </summary>
        public InputList<object> Triggers
        {
            get => _triggers ?? (_triggers = new InputList<object>());
            set => _triggers = value;
        }

        /// <summary>
        /// The command to run on update, if empty, create will run again.
        /// </summary>
        [Input("update")]
        public Input<string>? Update { get; set; }

        public CommandArgs()
        {
        }
        public static new CommandArgs Empty => new CommandArgs();
    }
}
