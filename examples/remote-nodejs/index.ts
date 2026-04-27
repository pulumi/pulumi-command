import { Config, asset, interpolate } from "@pulumi/pulumi";
import { local, remote, types } from "@pulumi/command";

const config = new Config();
const host = config.require("host");
const port = config.requireNumber("port");
const user = config.require("user");
const privateKey = Buffer.from(config.require("privateKeyBase64"), "base64").toString("ascii");

const connection: types.input.remote.ConnectionArgs = {
    host,
    port,
    user,
    privateKey,
};

const connectionNoDialRetry: types.input.remote.ConnectionArgs = {
    ...connection,
    dialErrorLimit: 1,
};

// Poll until the server accepts a connection. Other commands depend on this so
// they run only once the SSH server is reachable.
const poll = new remote.Command("poll", {
    connection: { ...connection, dialErrorLimit: -1 },
    create: "echo 'Connection established'",
}, { customTimeouts: { create: "2m" } });

const hostname = new remote.Command("hostname", {
    connection,
    create: "hostname",
}, { dependsOn: poll });

new remote.Command("remoteWrite", {
    connection,
    create: `echo hello > /tmp/remote-nodejs.txt`,
    delete: `rm /tmp/remote-nodejs.txt`,
}, { deleteBeforeReplace: true, dependsOn: poll });

new remote.Command("remoteWithNoDialRetry", {
    connection: connectionNoDialRetry,
    create: `echo hello > /tmp/remote-nodejs_no_dial_retry.txt`,
    delete: `rm /tmp/remote-nodejs_no_dial_retry.txt`,
}, { deleteBeforeReplace: true, dependsOn: poll });

new local.Command("localWrite", {
    create: `echo hello > local_only.txt`,
    delete: `rm local_only.txt`,
}, { deleteBeforeReplace: true });

const sizeFile = new remote.CopyToRemote("size", {
    connection,
    source: new asset.StringAsset("micro\n"),
    remotePath: "/tmp/size.txt",
}, { dependsOn: poll });

const catSize = new remote.Command("checkSize", {
    connection,
    create: "cat /tmp/size.txt",
}, { dependsOn: sizeFile });

export const connectionSecret = hostname.connection;
export const confirmSize = catSize.stdout;
export const hostnameStdout = hostname.stdout;
