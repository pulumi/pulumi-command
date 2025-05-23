// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../types/input";
import * as outputs from "../types/output";
import * as enums from "../types/enums";

import * as utilities from "../utilities";

export namespace remote {
    /**
     * Instructions for how to connect to a remote endpoint.
     */
    export interface Connection {
        /**
         * SSH Agent socket path. Default to environment variable SSH_AUTH_SOCK if present.
         */
        agentSocketPath?: string;
        /**
         * Max allowed errors on trying to dial the remote host. -1 set count to unlimited. Default value is 10.
         */
        dialErrorLimit?: number;
        /**
         * The address of the resource to connect to.
         */
        host: string;
        /**
         * The expected host key to verify the server's identity. If not provided, the host key will be ignored.
         */
        hostKey?: string;
        /**
         * The password we should use for the connection.
         */
        password?: string;
        /**
         * Max number of seconds for each dial attempt. 0 implies no maximum. Default value is 15 seconds.
         */
        perDialTimeout?: number;
        /**
         * The port to connect to. Defaults to 22.
         */
        port?: number;
        /**
         * The contents of an SSH key to use for the connection. This takes preference over the password if provided.
         */
        privateKey?: string;
        /**
         * The password to use in case the private key is encrypted.
         */
        privateKeyPassword?: string;
        /**
         * The connection settings for the bastion/proxy host.
         */
        proxy?: outputs.remote.ProxyConnection;
        /**
         * The user that we should use for the connection.
         */
        user?: string;
    }
    /**
     * connectionProvideDefaults sets the appropriate defaults for Connection
     */
    export function connectionProvideDefaults(val: Connection): Connection {
        return {
            ...val,
            dialErrorLimit: (val.dialErrorLimit) ?? 10,
            perDialTimeout: (val.perDialTimeout) ?? 15,
            port: (val.port) ?? 22,
            proxy: (val.proxy ? outputs.remote.proxyConnectionProvideDefaults(val.proxy) : undefined),
            user: (val.user) ?? "root",
        };
    }

    /**
     * Instructions for how to connect to a remote endpoint via a bastion host.
     */
    export interface ProxyConnection {
        /**
         * SSH Agent socket path. Default to environment variable SSH_AUTH_SOCK if present.
         */
        agentSocketPath?: string;
        /**
         * Max allowed errors on trying to dial the remote host. -1 set count to unlimited. Default value is 10.
         */
        dialErrorLimit?: number;
        /**
         * The address of the bastion host to connect to.
         */
        host: string;
        /**
         * The expected host key to verify the server's identity. If not provided, the host key will be ignored.
         */
        hostKey?: string;
        /**
         * The password we should use for the connection to the bastion host.
         */
        password?: string;
        /**
         * Max number of seconds for each dial attempt. 0 implies no maximum. Default value is 15 seconds.
         */
        perDialTimeout?: number;
        /**
         * The port of the bastion host to connect to.
         */
        port?: number;
        /**
         * The contents of an SSH key to use for the connection. This takes preference over the password if provided.
         */
        privateKey?: string;
        /**
         * The password to use in case the private key is encrypted.
         */
        privateKeyPassword?: string;
        /**
         * The user that we should use for the connection to the bastion host.
         */
        user?: string;
    }
    /**
     * proxyConnectionProvideDefaults sets the appropriate defaults for ProxyConnection
     */
    export function proxyConnectionProvideDefaults(val: ProxyConnection): ProxyConnection {
        return {
            ...val,
            dialErrorLimit: (val.dialErrorLimit) ?? 10,
            perDialTimeout: (val.perDialTimeout) ?? 15,
            port: (val.port) ?? 22,
            user: (val.user) ?? "root",
        };
    }

}
