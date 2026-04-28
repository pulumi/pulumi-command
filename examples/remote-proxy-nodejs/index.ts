import { Config, asset } from "@pulumi/pulumi";
import { remote, types } from "@pulumi/command";

const config = new Config();
const proxyHost = config.require("proxyHost");
const proxyPort = config.requireNumber("proxyPort");
const targetHost = config.require("targetHost");
const targetPort = config.requireNumber("targetPort");
const user = config.require("user");
const privateKey = Buffer.from(config.require("privateKeyBase64"), "base64").toString("ascii");

const proxyConnection: types.input.remote.ProxyConnectionArgs = {
    host: proxyHost,
    port: proxyPort,
    user,
    privateKey,
};

// targetHost is only resolvable from inside the Docker network (it's the proxy
// container's view of the second container), so the connection must go through
// the proxy.
const connection: types.input.remote.ConnectionArgs = {
    host: targetHost,
    port: targetPort,
    user,
    privateKey,
    proxy: proxyConnection,
};

const hostname = new remote.Command("hostname", {
    connection: {
        ...connection,
        dialErrorLimit: -1,
        proxy: { ...proxyConnection, dialErrorLimit: -1 },
    },
    create: "hostname",
}, { customTimeouts: { create: "2m" } });

new remote.Command("remoteWrite", {
    connection,
    create: `echo hello > /tmp/remote-proxy-nodejs.txt`,
    delete: `rm /tmp/remote-proxy-nodejs.txt`,
}, { deleteBeforeReplace: true, dependsOn: hostname });

const sizeFile = new remote.CopyToRemote("size", {
    connection,
    source: new asset.StringAsset("micro\n"),
    remotePath: "/tmp/size.txt",
}, { dependsOn: hostname });

const catSize = new remote.Command("checkSize", {
    connection,
    create: "cat /tmp/size.txt",
}, { dependsOn: sizeFile });

export const connectionSecret = hostname.connection;
export const confirmSize = catSize.stdout;
export const hostnameStdout = hostname.stdout;
