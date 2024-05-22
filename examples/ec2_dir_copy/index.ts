import * as pulumi from "@pulumi/pulumi";
import { remote, types } from "@pulumi/command";
import * as aws from "@pulumi/aws";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import { hashElement } from "folder-hash";
import { size } from "./size";

export = async () => {
    const config = new pulumi.Config();
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
            values: ["al2023-ami-2023.*-kernel-*-x86_64"],
        }],
    });

    const server = new aws.ec2.Instance("server", {
        instanceType: size,
        ami: ami.id,
        keyName: keyName,
        vpcSecurityGroupIds: [secgrp.id],
    }, {
        replaceOnChanges: ["instanceType"],
    });

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

    const from = path.join(__dirname, "src/");
    const to = config.get("destDir")!;

    const archive = new pulumi.asset.FileArchive(from);
    const copy = new remote.Copy("copy", {
        connection,
        archive: archive,
        remotePath: to,
    }, { dependsOn: poll });

    // Run `ls` on the remote to verify that the expected files were copied there.
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