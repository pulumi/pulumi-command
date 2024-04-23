import { interpolate, Config, secret } from "@pulumi/pulumi";
import { local, remote, types } from "@pulumi/command";
import * as aws from "@pulumi/aws";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

const size = "t2.nano";

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
}, { customTimeouts: { create: "10m" } })

////// Start of the actual test //////

// var from = path.join(os.tmpdir(), "from");
// fs.closeSync(fs.openSync(path.join(from, "file1"), 'w'));
// fs.mkdirSync(path.join(from, "dir1"));
// fs.closeSync(fs.openSync(path.join(from, "dir1", "file2"), 'w'));

// const to = "/tmp/to";

const from = config.get("srcDir")!;
const to = config.get("destDir")!;

const copy = new remote.CopyFile("copy", {
    connection,
    localPath: from,
    remotePath: to,
})

const ls = new remote.Command("ls", {
    connection,
    create: `find ${to} | sort`, // -R is recursive, -p shows directories via trailing slash
}, { dependsOn: copy });

export const destDir = to;
export const lsRemote = ls.stdout;
