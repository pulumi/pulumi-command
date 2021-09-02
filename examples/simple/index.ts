import * as command from "@pulumi/command";

const random = new command.Random("my-random", { length: 24 });

export const output = random.result;