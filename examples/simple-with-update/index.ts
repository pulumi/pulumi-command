import * as local from "@pulumi/command/local";
import * as random from "@pulumi/random";
import { interpolate } from "@pulumi/pulumi";
import { len, fail, update } from "./extras";

const pw = new random.RandomPassword("pw", { length: len, special: false });

const pwd = new local.Command("pwd", {
    create: interpolate`touch "${pw.result}cat.txt"`,
    delete: interpolate`rm "${pw.result}dog.txt"`,
    triggers: [pw.result],
})

export const output = pwd.stdout;
