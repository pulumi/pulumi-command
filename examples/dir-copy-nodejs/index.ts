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

    // Exercise the remaining asset/archive subtypes that CopyToRemote supports. Each is written to
    // its own subdirectory under `extrasTo` so the verification command below can list them all.
    const extrasTo = `${to}-extras`;

    const stringAssetCopy = new remote.CopyToRemote("string-asset", {
        connection,
        source: new pulumi.asset.StringAsset("hello from a string asset\n"),
        remotePath: `${extrasTo}/string-asset.txt`,
    }, { dependsOn: poll });

    const remoteAssetCopy = new remote.CopyToRemote("remote-asset", {
        connection,
        // file:// URIs are resolved on the machine running Pulumi, so this exercises the
        // RemoteAsset code path without depending on a public HTTP server.
        source: new pulumi.asset.RemoteAsset("file://" + path.join(__dirname, "src/file1")),
        remotePath: `${extrasTo}/remote-asset.txt`,
    }, { dependsOn: poll });

    const remoteArchiveCopy = new remote.CopyToRemote("remote-archive", {
        connection,
        source: new pulumi.asset.RemoteArchive("file://" + path.join(__dirname, "fixtures/sample.tar.gz")),
        // Remote archives are copied as-is, so the destination is a file path.
        remotePath: `${extrasTo}/remote-archive.tar.gz`,
    }, { dependsOn: poll });

    const assetArchiveCopy = new remote.CopyToRemote("asset-archive", {
        connection,
        source: new pulumi.asset.AssetArchive({
            "greeting.txt": new pulumi.asset.StringAsset("hello\n"),
            "nested/answer.txt": new pulumi.asset.StringAsset("42\n"),
        }),
        remotePath: `${extrasTo}/asset-archive`,
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

    const lsExtras = new remote.Command("ls-extras", {
        connection,
        create: `find ${extrasTo} | sort`,
        triggers: [hash],
    }, { dependsOn: [stringAssetCopy, remoteAssetCopy, remoteArchiveCopy, assetArchiveCopy] });

    return {
        destDir: to,
        lsRemote: ls.stdout,
        lsExtras: lsExtras.stdout,
    };
};
