import * as local from "@pulumi/command/local";
import * as random from "@pulumi/random";
import { interpolate } from "@pulumi/pulumi";

const pw = new random.RandomPassword("pw", { length: 10 });

const plainFile = local.runOutput({
    command: `echo "Hello world!" > hello.txt`,
    assets: ["hello.txt"]
});

const secretFile = local.runOutput({
    command: interpolate`echo "${pw.result}" > password.txt`,
    assets: ["password.txt"]
});

export const plainOutput = plainFile.stdout;
export const plainAssets = plainFile.assets;
export const secretOutput = secretFile.stdout;
export const secretAssets = secretFile.assets;
