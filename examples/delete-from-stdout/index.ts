import { local } from "@pulumi/command";

const mktemp = new local.Command('mktemp', {
    create: 'mktemp',
    delete: 'rm $PULUMI_COMMAND_STDOUT'
})

export const output = mktemp.stdout;
