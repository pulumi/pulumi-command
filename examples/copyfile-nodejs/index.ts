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

    const from = path.join(__dirname, "src/file1");
    const copy = new remote.CopyFile("copy", {
        connection,
        localPath: from,
        remotePath: to,
    }, { dependsOn: poll });

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
