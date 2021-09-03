import * as command from "@pulumi/command";
import * as random from "@pulumi/random";
import * as aws from "@pulumi/aws";
import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config();
const publicKey = config.require("publicKey");
const privateKey = config.requireSecret("privateKey").apply(key => {
    if (key.startsWith("-----BEGIN RSA PRIVATE KEY-----")) {
        return key;
    } else {
        return Buffer.from(key, "base64").toString("ascii");
    }
});

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
const size = "t2.micro";
const key = new aws.ec2.KeyPair("key", { publicKey });
const server = new aws.ec2.Instance("server", {
    instanceType: size,
    ami: amiId,
    keyName: key.keyName,
    vpcSecurityGroupIds: [secgrp.id],
});

const user = "ec2-user";
const host = server.publicIp;

const cpConfig = new command.Command("config", {
    create: pulumi.interpolate`
    
        ssh-add - <<< "${privateKey}" 2>/dev/null
        scp myapp.conf ${user}@${host}:myapp.conf
        `,
})

const catConfig = new command.Command("cat", {
    create: pulumi.interpolate`

        ssh-add - <<< "${privateKey}" 2>/dev/null
        ssh ${user}@${host} 'cat myapp.conf'
        `,
}, { dependsOn: [cpConfig] });

export const publicIp = server.publicIp;
export const publicHostName = server.publicDns;
export const catConfigStdout = catConfig.stdout;