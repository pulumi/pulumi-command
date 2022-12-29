import * as local from "@pulumi/command/local";
import * as random from "@pulumi/random";
import { interpolate } from "@pulumi/pulumi";
import { len, fail, update } from "./extras";

const pw = new random.RandomPassword("pw", { length: len, special: false });

const pwd = new local.Command("pwd", {
    create: interpolate`echo "${pw.result}" > password.txt`,
    delete: `rm -f password.txt`,
    triggers: [pw.result],
}, { deleteBeforeReplace: true });


let deleteCommand = "rm -f password2.txt";
if (update) {
    deleteCommand += " && echo 'deleted'";
}
const pwd2 = new local.Command("pwd2", {
    create: `echo "$PASSWORD" > password2.txt`,
    delete: deleteCommand,
    environment: {
        PASSWORD: pw.result,
    },
    triggers: [pw.result],
}, { deleteBeforeReplace: true });

// Manage an external artifact which is created after a resource is created, 
// deleted before a resource is destroyed, and recreated when the resource is 
// replaced.  This could also register/deregister the resource with an external 
// registration or other remote API instead of just writing to local disk.
const pwd3 = new local.Command("pwd3", {
    create: interpolate`touch "${pw.result}.txt"`,
    delete: interpolate`rm "${pw.result}.txt"`,
    triggers: [pw.result],
})

if (fail) {
    new local.Command("fail", {
        create: `echo "couldn't do what I wanted..." && false`,
    });
}

export const output = pwd.stdout;
