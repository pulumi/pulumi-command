import { interpolate } from "@pulumi/pulumi";
import * as command from "@pulumi/command";
import * as aws from "@pulumi/aws";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

const publicKey = fs.readFileSync(path.join(os.homedir(), ".ssh", "id_rsa.pub")).toString("utf8");
const privateKey = fs.readFileSync(path.join(os.homedir(), ".ssh", "id_rsa")).toString("utf8");

// Create a new security group that permits SSH and web access.
const secgrp = new aws.ec2.SecurityGroup("secgrp", {
    description: "Foo",
    ingress: [
        { protocol: "tcp", fromPort: 22, toPort: 22, cidrBlocks: ["0.0.0.0/0"] },
        { protocol: "tcp", fromPort: 80, toPort: 80, cidrBlocks: ["0.0.0.0/0"] },
    ],
});

// Get the AMI
const amiId = aws.ec2.getAmi({
    owners: ["amazon"],
    mostRecent: true,
    filters: [{
        name: "name",
        values: ["amzn2-ami-hvm-2.0.????????-x86_64-gp2"],
    }],
}, { async: true }).then(ami => ami.id);

// Create an EC2 server that we'll then provision stuff onto.
const size = "t2.nano";
const key = new aws.ec2.KeyPair("key", { publicKey });
const server = new aws.ec2.Instance("server", {
    instanceType: size,
    ami: amiId,
    keyName: key.keyName,
    vpcSecurityGroupIds: [secgrp.id],
}, { replaceOnChanges: ["instanceType"] });

const connection: command.types.input.RemoteConnectionArgs = {
    host: server.publicIp,
    user: "ec2-user",
    privateKey: privateKey,
};

const hostname = new command.RemoteCommand("hostname", {
    connection,
    create: "hostname",
});

const remotePrivateIP = new command.RemoteCommand("remotePrivateIP", {
    connection,
    create: interpolate`echo ${server.privateIp} > private_ip.txt`,
    delete: `rm private_ip.txt`,
}, { deleteBeforeReplace: true });


const localPrivateIP = new command.Command("localPrivateIP", {
    create: interpolate`echo ${server.privateIp} > private_ip.txt`,
    delete: `rm private_ip.txt`,
}, { deleteBeforeReplace: true });

export const publicIp = server.publicIp;
export const publicHostName = server.publicDns;
export const hostnameStdout = hostname.stdout;
