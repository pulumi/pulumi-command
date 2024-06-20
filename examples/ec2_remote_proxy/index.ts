import { interpolate, Config, asset } from "@pulumi/pulumi";
import { local, remote, types } from "@pulumi/command";
import * as aws from "@pulumi/aws";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import { size } from "./size";

const config = new Config();
const keyName = config.get("keyName") ?? new aws.ec2.KeyPair("key", { publicKey: config.require("publicKey") }).keyName;
const privateKeyBase64 = config.get("privateKeyBase64");
const privateKey = privateKeyBase64 ? Buffer.from(privateKeyBase64, 'base64').toString('ascii') : fs.readFileSync(path.join(os.homedir(), ".ssh", "id_rsa")).toString("utf8");

const ingress = new aws.ec2.SecurityGroup("ingress", {
    description: "A security group that will accept SSH connections from the outside world.",
    ingress: [
        { protocol: "tcp", fromPort: 22, toPort: 22, cidrBlocks: ["0.0.0.0/0"] },
        { protocol: "tcp", fromPort: 80, toPort: 80, cidrBlocks: ["0.0.0.0/0"] },
    ],
    egress: [
        { fromPort: 0, toPort: 0, protocol: "-1", cidrBlocks: ["0.0.0.0/0"], ipv6CidrBlocks: ["::/0"], },
    ],
});

const validated = new aws.ec2.SecurityGroup("validated", {
    description: "A security group that will only accept connections that have already been validated.",
    ingress: [
        { protocol: "tcp", fromPort: 22, toPort: 22, securityGroups: [ingress.id] },
        { protocol: "tcp", fromPort: 80, toPort: 80, securityGroups: [ingress.id] },
    ],
    egress: [
        { fromPort: 0, toPort: 0, protocol: "-1", cidrBlocks: ["0.0.0.0/0"], ipv6CidrBlocks: ["::/0"], },
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
    vpcSecurityGroupIds: [validated.id],
}, { replaceOnChanges: ["instanceType"] });

const proxy = new aws.ec2.Instance("proxy", {
    instanceType: size,
    ami: ami.id,
    keyName: keyName,
    vpcSecurityGroupIds: [ingress.id],
}, { replaceOnChanges: ["instanceType"] });

const connection: types.input.remote.ConnectionArgs = {
    host: server.privateDns,
    user: "ec2-user",
    privateKey: privateKey,
    proxy: {
        host: proxy.publicIp,
        user: "ec2-user",
        privateKey: privateKey,
    },
};

const hostname = new remote.Command("hostname", {
    connection: {
        ...connection,
        dialErrorLimit: -1,
        proxy: {
            ...connection.proxy,
            dialErrorLimit: -1,
        }
    },
    create: "hostname",
}, { customTimeouts: { create: "12m" } });

new remote.Command("remotePrivateIP", {
    connection,
    create: interpolate`echo ${server.privateIp} > private_ip.txt`,
    delete: `rm private_ip.txt`,
}, { deleteBeforeReplace: true, dependsOn: hostname });

new local.Command("localPrivateIP", {
    create: interpolate`echo ${server.privateIp} > private_ip.txt`,
    delete: `rm private_ip.txt`,
}, { deleteBeforeReplace: true, dependsOn: hostname });

const sizeFile = new remote.CopyToRemote("size", {
    connection: connection,
    source: new asset.FileAsset("./size.ts"),
    remotePath: "size.ts",
}, { dependsOn: hostname })

const catSize = new remote.Command("checkSize", {
    connection: connection,
    create: "cat size.ts",
}, { dependsOn: sizeFile })

export const connectionSecret = hostname.connection;
export const confirmSize = catSize.stdout;
export const publicIp = server.publicIp;
export const publicHostName = server.publicDns;
export const hostnameStdout = hostname.stdout;
