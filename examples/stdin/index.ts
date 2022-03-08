import { local } from "@pulumi/command";

const random = new local.Command("stdin", {
    create: "head -n 1",
    stdin: "the quick brown fox\njumped over\nthe lazy dog"
});

export const output = random.stdout;
