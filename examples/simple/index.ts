import * as command from "@pulumi/command";
import * as random from "@pulumi/random";
import { interpolate } from "@pulumi/pulumi";

const pw = new random.RandomPassword("pw", { length: 10 });

const pwd = new command.Command("pwd", {
    create: interpolate`echo ${pw.result} > password.txt`,
    delete: "rm password.txt",
}, { ignoreChanges: ["create"] });

export const output = pwd.stdout;