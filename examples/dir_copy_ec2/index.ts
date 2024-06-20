import * as pulumi from "@pulumi/pulumi";
import { remote, types } from "@pulumi/command";
import * as aws from "@pulumi/aws-native";
import * as path from "path";
import { hashElement } from "folder-hash";
import { size } from "./size";

export = async () => {
    const config = new pulumi.Config();
    const keyName = "sshKey"

    const keyPair = new aws.ec2.KeyPair("keyPair", {
        keyName: keyName,
    });

    const secgrp = new aws.ec2.SecurityGroup("secgrp", {
        groupDescription: "Foo",
        securityGroupIngress: [
            { ipProtocol: "tcp", fromPort: 22, toPort: 22, cidrIp: "0.0.0.0/0" },
            // { protocol: "tcp", fromPort: 80, toPort: 80, cidrBlocks: ["0.0.0.0/0"] },
        ],
    });

    const server = new aws.ec2.Instance("server", {
        instanceType: size,
        imageId: "ami-0e58172bedd62916b", //
        keyName: keyName,
        securityGroupIds: [secgrp.id],
    }, {
        replaceOnChanges: ["instanceType"],
    });

    var privateKey:pulumi.Output<string>
    const getPrivateKey = () => {
        if (privateKey) {
            return privateKey
        }
        privateKey = aws.ssm.getParameterOutput({ name: pulumi.interpolate`"/ec2/keypair/${keyPair.keyPairId}` }).value!.apply(x => x!)
        return privateKey
    }

    const connection: types.input.remote.ConnectionArgs = {
        host: server.publicIp,
        user: "ec2-user",
    };

    // We poll the server until it responds.
    //
    // Because other commands depend on this command, other commands are guaranteed
    // to hit an already booted server.
    const poll = new remote.Command("poll", {
        connection: { ...connection, dialErrorLimit: -1, privateKey: getPrivateKey() },
        create: "echo 'Connection established'",
    }, {
        customTimeouts: { create: "10m" },
        dependsOn: keyPair,
    })

    ////// Start of the actual test //////

    const from = path.join(__dirname, "src/");
    const to = config.get("destDir")!;

    const archive = new pulumi.asset.FileArchive(from);
    const copy = new remote.Copy("copy", {
        connection: { ...connection, privateKey: getPrivateKey() },
        source: archive,
        remotePath: to,
    }, { dependsOn: poll });

    // Verify that the expected files were copied to the remote.
    // We want to run this after each copy, i.e., when something changed, but not otherwise to avoid unclean refreshes.
    // We use the hash of the source directory as a trigger to achieve this, since the trigger needs to be a primitive
    // value and we cannot use the Copy resource itself.
    // The FileArchive already a hash calculated, but it's not exposed.
    const hash = await hashElement(from);
    const ls = new remote.Command("ls", {
        connection: { ...connection, privateKey: getPrivateKey() },
        create: `find ${to} | sort`,
        triggers: [hash],
    }, { dependsOn: copy });

    return {
        destDir: to,
        lsRemote: ls.stdout
    }
}