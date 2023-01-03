import * as command from "@pulumi/command";

new command.local.Command("stdout-and-stderr-success", {
  create: "ls not-a-file index.ts not-a-file-2 || true"
});

new command.local.Command("stdout-and-stderr-error", {
    create: "ls not-a-file index.ts not-a-file-2"
  });