import * as pulumi from "@pulumi/pulumi";
import * as random from "@pulumi/random";
import { local } from "@pulumi/command";

interface LabelArgs {
    owner: pulumi.Input<string>;
    repo: pulumi.Input<string>;
    name: pulumi.Input<string>;
    githubToken: pulumi.Input<string>;
}

class GitHubLabel extends pulumi.ComponentResource {
    public url: pulumi.Output<string>;

    constructor(name: string, args: LabelArgs, opts?: pulumi.ComponentResourceOptions) {
        super("example:github:Label", name, args, opts);

        const label = new local.Command("label", {
            create: "./create_label.sh",
            delete: "./delete_label.sh",
            environment: {
                OWNER: args.owner,
                REPO: args.repo,
                NAME: args.name,
                GITHUB_TOKEN: args.githubToken,
            }
        }, { parent: this });

        const response = label.stdout.apply(JSON.parse);
        this.url = response.apply((x: any) => x.url as string);
    }
}

const config = new pulumi.Config();
const rand = new random.RandomString("s", { length: 10, special: false });

const label = new GitHubLabel("l", {
    owner: "pulumi",
    repo: "pulumi-command",
    name: rand.result,
    githubToken: config.requireSecret("githubToken"),
});

export const labelUrl = label.url;
