import * as local from "@pulumi/command/local";
import * as random from "@pulumi/random";
import { interpolate } from "@pulumi/pulumi";
import { len, fail } from "./extras";

const pw = new random.RandomPassword("pw", { length: len });

const pwd = new local.Command("pwd", {
    create: interpolate`echo "${pw.result}" > password.txt`,
    delete: `rm -f password.txt`,
    replaceOnChanges: [pw.result],
}, { deleteBeforeReplace: true });

const pwd2 = new local.Command("pwd2", {
    create: `echo "$PASSWORD" > password2.txt`,
    delete: `rm -f password2.txt`,
    environment: {
        PASSWORD: pw.result,
    },
    replaceOnChanges: [pw.result],
}, { deleteBeforeReplace: true });

// Manage an external artifact which is created after a resource is created, 
// deleted before a resource is destroyed, and recreated when the resource is 
// replaced.  This could also register/deregister the resource with an external 
// registration or otehr remote API instead of just writing to local disk.
const pwd3 = new local.Command("pwd3", {
    create: interpolate`touch "${pw.result}.txt"`,
    delete: interpolate`rm "${pw.result}.txt"`,
    replaceOnChanges: [pw.result],
})

if (fail) {
    new local.Command("fail", {
        create: `echo "couldn't do what I wanted..." && false`,
    });
}

export const output = pwd.stdout;
