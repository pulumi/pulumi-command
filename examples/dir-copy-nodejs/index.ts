import * as pulumi from "@pulumi/pulumi";
import { remote, types } from "@pulumi/command";
import * as path from "path";
import { hashElement } from "folder-hash";

export = async () => {
    const config = new pulumi.Config();
    const host = config.require("host");
    const port = config.requireNumber("port");
    const user = config.require("user");
    const privateKey = Buffer.from(config.require("privateKeyBase64"), "base64").toString("ascii");
    const to = config.require("destDir");

    const connection: types.input.remote.ConnectionArgs = {
        host,
        port,
        user,
        privateKey,
    };

    const poll = new remote.Command("poll", {
        connection: { ...connection, dialErrorLimit: -1 },
        create: "echo 'Connection established'",
    }, { customTimeouts: { create: "2m" } });

    const from = path.join(__dirname, "src/");
    const archive = new pulumi.asset.FileArchive(from);
    const copy = new remote.CopyToRemote("copy", {
        connection,
        source: archive,
        remotePath: to,
    }, { dependsOn: poll });

    // Use the source-tree hash as a trigger so `ls` only re-runs when the
    // copy itself changed. FileArchive computes a hash internally but doesn't
    // expose it.
    const hash = await hashElement(from);
    const ls = new remote.Command("ls", {
        connection,
        create: `find ${to} | sort`,
        triggers: [hash],
    }, { dependsOn: copy });

    return {
        destDir: to,
        lsRemote: ls.stdout,
    };
};
