import { interpolate, Config, asset } from "@pulumi/pulumi";
import { local, remote, types } from "@pulumi/command";
import * as aws from "@pulumi/aws";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import { size } from "./size";

const config = new Config();
const keyName = config.get("keyName") ??
  new aws.ec2.KeyPair("key", { publicKey: config.require("publicKey") }).keyName;
const privateKeyBase64 = config.get("privateKeyBase64");
const privateKey = privateKeyBase64 ?
  Buffer.from(privateKeyBase64, 'base64').toString('ascii') :
  fs.readFileSync(path.join(os.homedir(), ".ssh", "id_rsa")).toString("utf8");

const secgrp = new aws.ec2.SecurityGroup("secgrp", {
    description: "Foo",
    ingress: [
        { protocol: "tcp", fromPort: 22, toPort: 22, cidrBlocks: ["0.0.0.0/0"] },
        { protocol: "tcp", fromPort: 80, toPort: 80, cidrBlocks: ["0.0.0.0/0"] },
    ],
});

const ami = aws.ec2.getAmiOutput({
    owners: ["amazon"],
    mostRecent: true,
    filters: [{
        name: "name",
        values: ["amzn2-ami-hvm-*-x86_64-ebs"],
    }],
});

const server = new aws.ec2.Instance("server", {
    instanceType: size,
    ami: ami.id,
    keyName: keyName,
    vpcSecurityGroupIds: [secgrp.id],
}, { replaceOnChanges: ["instanceType"] });

const connection: types.input.remote.ConnectionArgs = {
    host: server.publicIp,
    user: "ec2-user",
    privateKey: privateKey,
};

const connectionNoDialRetry: types.input.remote.ConnectionArgs = {
    ...connection,
    dialErrorLimit: 1,
};

// We poll the server until it responds.
//
// Because other commands depend on this command, other commands are guaranteed
// to hit an already booted server.
const poll = new remote.Command("poll", {
  connection: { ...connection, dialErrorLimit: -1 },
    create: "echo 'Connection established'",
}, { customTimeouts: { create: "12m" } })

const hostname = new remote.Command("hostname", {
    connection,
    create: "hostname",
}, { dependsOn: poll });

new remote.Command("remotePrivateIP", {
    connection,
    create: interpolate`echo ${server.privateIp} > private_ip.txt`,
    delete: `rm private_ip.txt`,
}, { deleteBeforeReplace: true, dependsOn: poll });

new remote.Command("remoteWithNoDialRetryPrivateIP", {
    connection: connectionNoDialRetry,
    create: interpolate`echo ${server.privateIp} > private_ip_on_no_dial_retry.txt`,
    delete: `rm private_ip_on_no_dial_retry.txt`,
}, { deleteBeforeReplace: true, dependsOn: poll });

new local.Command("localPrivateIP", {
    create: interpolate`echo ${server.privateIp} > private_ip.txt`,
    delete: `rm private_ip.txt`,
}, { deleteBeforeReplace: true });

const sizeFile = new remote.CopyToRemote("size", {
    connection,
    source: new asset.FileAsset("./size.ts"),
    remotePath: "size.ts",
}, { dependsOn: poll })

const catSize = new remote.Command("checkSize", {
    connection,
    create: "cat size.ts",
}, { dependsOn: sizeFile })

export const connectionSecret = hostname.connection;
export const confirmSize = catSize.stdout;
export const publicIp = server.publicIp;
export const publicHostName = server.publicDns;
export const hostnameStdout = hostname.stdout;
