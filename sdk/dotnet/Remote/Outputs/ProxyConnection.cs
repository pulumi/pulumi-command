// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Command.Remote.Outputs
{

    /// <summary>
    /// Instructions for how to connect to a remote endpoint via a bastion host.
    /// </summary>
    [OutputType]
    public sealed class ProxyConnection
    {
        /// <summary>
        /// SSH Agent socket path. Default to environment variable SSH_AUTH_SOCK if present.
        /// </summary>
        public readonly string? AgentSocketPath;
        /// <summary>
        /// Max allowed errors on trying to dial the remote host. -1 set count to unlimited. Default value is 10.
        /// </summary>
        public readonly int? DialErrorLimit;
        /// <summary>
        /// The address of the bastion host to connect to.
        /// </summary>
        public readonly string Host;
        /// <summary>
        /// The password we should use for the connection to the bastion host.
        /// </summary>
        public readonly string? Password;
        /// <summary>
        /// Max number of seconds for each dial attempt. 0 implies no maximum. Default value is 15 seconds.
        /// </summary>
        public readonly int? PerDialTimeout;
        /// <summary>
        /// The port of the bastion host to connect to.
        /// </summary>
        public readonly double? Port;
        /// <summary>
        /// The contents of an SSH key to use for the connection. This takes preference over the password if provided.
        /// </summary>
        public readonly string? PrivateKey;
        /// <summary>
        /// The password to use in case the private key is encrypted.
        /// </summary>
        public readonly string? PrivateKeyPassword;
        /// <summary>
        /// The user that we should use for the connection to the bastion host.
        /// </summary>
        public readonly string? User;

        [OutputConstructor]
        private ProxyConnection(
            string? agentSocketPath,

            int? dialErrorLimit,

            string host,

            string? password,

            int? perDialTimeout,

            double? port,

            string? privateKey,

            string? privateKeyPassword,

            string? user)
        {
            AgentSocketPath = agentSocketPath;
            DialErrorLimit = dialErrorLimit;
            Host = host;
            Password = password;
            PerDialTimeout = perDialTimeout;
            Port = port;
            PrivateKey = privateKey;
            PrivateKeyPassword = privateKeyPassword;
            User = user;
        }
    }
}
