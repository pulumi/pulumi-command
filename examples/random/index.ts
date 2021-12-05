import * as command from "@pulumi/command";

const random = new command.Command("random", {
    create: "openssl rand -hex 16",
})

export const output = random.stdout;