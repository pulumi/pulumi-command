import { asset, interpolate, Config } from "@pulumi/pulumi";
import { local, remote, types } from "@pulumi/command";
import * as aws from "@pulumi/aws";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

const config = new Config();
const keyName = config.get("keyName") ?? new aws.ec2.KeyPair("key", { publicKey: config.require("publicKey") }).keyName;
const privateKeyBase64 = config.get("privateKeyBase64");
const privateKey = privateKeyBase64 ? Buffer.from(privateKeyBase64, 'base64').toString('ascii') : fs.readFileSync(path.join(os.homedir(), ".ssh", "id_rsa")).toString("utf8");

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
        values: ["amzn2-ami-hvm-2.0.????????-x86_64-gp2"],
    }],
});

const server = new aws.ec2.Instance("server", {
    instanceType: "t2.nano",
    ami: ami.id,
    keyName: keyName,
    vpcSecurityGroupIds: [secgrp.id],
}, { replaceOnChanges: ["instanceType"] });

const connection: types.input.remote.ConnectionArgs = {
    host: server.publicIp,
    user: "ec2-user",
    privateKey: privateKey,
};

const checkPython = new remote.Command("versionCheck", {
    connection,
    create: "python --version",
});

const program = new remote.CopyFile("program", {
    connection,
    localPath: new asset.FileAsset("./program.py"),
    remotePath: "program.py",
});

const zipDataPath = "/tmp/data.zip";

const zipData = new local.Command("zip", {
    create: `zip ${zipDataPath} data/*`,
    delete: `rm ${zipDataPath}`,
});

const data = new remote.CopyFile("data", {
    connection,
    localPath: zipData.stdout.apply((_: string) => new asset.FileArchive("/tmp/data.zip")),
    remotePath: "data",
});

// Depends on program and data via the interpolate.
const run = new remote.Command("run", {
    connection,
    create: interpolate`python ${program.remotePath} ${data.remotePath}/data`,
});


export const pythonVersion = checkPython.stdout.apply((x: string) => x.trim());
export const result = run.stdout;
