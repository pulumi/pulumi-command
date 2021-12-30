import { local } from "@pulumi/command";

const random = new local.Command("random", {
    create: "openssl rand -hex 16",
});

export const output = random.stdout;
