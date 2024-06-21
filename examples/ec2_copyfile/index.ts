import * as pulumi from "@pulumi/pulumi";
import { remote, types } from "@pulumi/command";
import * as aws from "@pulumi/aws";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import { hashElement } from "folder-hash";
import { size } from "./size";

// This test for remote.CopyFile duplicates a bunch of code from ec2_dir_copy, but since
// CopyFile is deprecated and will be removed, we don't need to spend time refactoring it.

export = async () => {
    // Get a key pair to connect to the EC2 instance. If the name of an existing key pair is
    // provided, use it, otherwise create one. We get the private key from config, or default to
    // the default id_rsa SSH key.
    const config = new pulumi.Config();
    const keyName = config.get("keyName") ??
        new aws.ec2.KeyPair("key", { publicKey: config.require("publicKey") }).keyName;
    const privateKeyBase64 = config.get("privateKeyBase64");
    const privateKey = privateKeyBase64 ?
        Buffer.from(privateKeyBase64, 'base64').toString('ascii') :
        fs.readFileSync(path.join(os.homedir(), ".ssh", "id_rsa")).toString("utf8");

    // Create a security group that allows SSH traffic.
    const secgrp = new aws.ec2.SecurityGroup("secgrp", {
        description: "Foo",
        ingress: [
            { protocol: "tcp", fromPort: 22, toPort: 22, cidrBlocks: ["0.0.0.0/0"] },
        ],
    });

    // Get the latest Amazon Linux AMI (image) for the region we're using.
    const ami = aws.ec2.getAmiOutput({
        owners: ["amazon"],
        mostRecent: true,
        filters: [{
            name: "name",
            values: ["al2023-ami-2023.*-kernel-*-x86_64"],
        }],
    });

    // Create the EC2 instance we will copy files to.
    const server = new aws.ec2.Instance("server", {
        instanceType: size,
        ami: ami.id,
        keyName: keyName,
        vpcSecurityGroupIds: [secgrp.id],
    }, {
        replaceOnChanges: ["instanceType"],
    });

    // The configuration of our SSH connection to the instance.
    const connection: types.input.remote.ConnectionArgs = {
        host: server.publicIp,
        user: "ec2-user",
        privateKey: privateKey,
    };

    // Poll the server until it responds.
    //
    // Because other commands depend on this command, other commands are guaranteed
    // to hit an already booted server.
    const poll = new remote.Command("poll", {
        connection: { ...connection, dialErrorLimit: -1 },
        create: "echo 'Connection established'",
    }, { customTimeouts: { create: "10m" } })

    ////// Start of the actual test //////

    const from = path.join(__dirname, "src/file1");
    const to = config.get("destDir")!;

    const copy = new remote.CopyFile("copy", {
        connection,
        localPath: from,
        remotePath: to,
    }, { dependsOn: poll });

    // Verify that the expected files were copied to the remote.
    // We want to run this after each copy, i.e., when something changed, but not otherwise to avoid unclean refreshes.
    // We use the hash of the source directory as a trigger to achieve this, since the trigger needs to be a primitive
    // value and we cannot use the Copy resource itself.
    // The FileArchive already a hash calculated, but it's not exposed.
    const hash = await hashElement(from);
    const ls = new remote.Command("ls", {
        connection,
        create: `find ${to} | sort`,
        triggers: [hash],
    }, { dependsOn: copy });

    return {
        destDir: to,
        lsRemote: ls.stdout
    }
}