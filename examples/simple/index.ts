import * as command from "@pulumi/command";
import * as random from "@pulumi/random";
import { interpolate } from "@pulumi/pulumi";

const pw = new random.RandomPassword("pw", { length: 10 });

const pwd = new command.Command("pwd", {
    create: interpolate`echo "${pw.result}" > password.txt`,
    delete: `rm -f password.txt`,
}, { deleteBeforeReplace: true });

const pwd2 = new command.Command("pwd2", {
    create: `echo "$PASSWORD" > password2.txt`,
    delete: `rm -f password2.txt`,
    environment: {
        PASSWORD: pw.result,
    }
}, { deleteBeforeReplace: true });

const fail = new command.Command("fail", {
    create: `echo "couldn't do what I wanted..." && false`,
});

export const output = pwd.stdout;