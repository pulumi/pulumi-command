import * as local from "@pulumi/command/local";
import * as random from "@pulumi/random";
import { interpolate } from "@pulumi/pulumi";

const pw = new random.RandomPassword("pw", { length: 10, special: false });

const plainFile = local.runOutput({
    command: `echo "Hello world!" > hello.txt`,
    assetPaths: ["*.txt", "!**password**"],
    archivePaths: ["*.txt", "!**password**"],
});

const secretFile = local.runOutput({
    command: interpolate`echo "${pw.result}" > password.txt`,
    assetPaths: ["password.txt"]
});

const globTest = local.runOutput({
    command: "pwd",
    dir: process.cwd(),
    archivePaths: [
        "**/*.txt",
        "*",
        "!yarn.lock",
        "!**password**",
    ]
})

export const plainOutput = plainFile.stdout;
export const plainAssets = plainFile.assets;
export const plainArchive = plainFile.archive;
export const secretOutput = secretFile.stdout;
export const secretAssets = secretFile.assets;
export const secretArchive = secretFile.archive;
export const globTestAssets = globTest.archive;
