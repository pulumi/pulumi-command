[![Actions Status](https://github.com/pulumi/pulumi-command/workflows/master/badge.svg)](https://github.com/pulumi/pulumi-command/actions)
[![Slack](http://www.pulumi.com/images/docs/badges/slack.svg)](https://slack.pulumi.com)
[![NPM version](https://badge.fury.io/js/%40pulumi%2Fcommand.svg)](https://www.npmjs.com/package/@pulumi/command)
[![Python version](https://badge.fury.io/py/pulumi-command.svg)](https://pypi.org/project/pulumi-command)
[![NuGet version](https://badge.fury.io/nu/pulumi.command.svg)](https://badge.fury.io/nu/pulumi.command)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/pulumi/pulumi-command/sdk/go)](https://pkg.go.dev/github.com/pulumi/pulumi-command/sdk/go)
[![License](https://img.shields.io/npm/l/%40pulumi%2Fpulumi.svg)](https://github.com/pulumi/pulumi-command/blob/master/LICENSE)

# Pulumi Command Provider (preview)

The Pulumi Command Provider enables you to execute commands and scripts either locally or remotely as part of the Pulumi resource model.  Resources in the command package support running scripts on `create` and `destroy` operations, supporting stateful local and remote command execution.

There are many scenarios where the Command package can be useful:

* Running a command locally after creating a resource, to register it with an external service
* Running a command locally before deleting a resource, to deregister it with an external service
* Running a command remotely on a remote host immediately after creating it
* Copying a file to a remote host after creating it (potentially as a script to be executed afterwards)
* As a simple alternative to some use cases for Dynamic Providers (especially in languages which do not yet support Dynamic Providers).

Some users may have experience with Terraform "provisioners", and the Command package offers support for similar scenarios.  However, the Command package is provided as independent resources which can be combined with other resources in many interesting ways. This has many strengths, but also some differences, such as the fact that a Command resource failing does not cause a resource it is operating on to fail.

You can use the Command package from a Pulumi program written in any Pulumi language: C#, Go, JavaScript/TypeScript, and Python.
You'll need to [install and configure the Pulumi CLI](https://pulumi.com/docs/get-started/install) if you haven't already.


> **NOTE**: The Command package is in preview.  The API design may change ahead of general availability based on [user feedback](https://github.com/pulumi/pulumi-command/issues). 

## Examples

### A simple local resource (random)

The simplest use case for `local.Command` is to just run a command on `create`, which can return some value which will be stored in the state file, and will be persistent for the life of the stack (or until the resource is destroyed or replaced).  The example below uses this as an alternative to the `random` package to create some randomness which is stored in Pulumi state.

```ts
import { local } from "@pulumi/command";

const random = new local.Command("random", {
    create: "openssl rand -hex 16",
});

export const output = random.stdout;
```

```go
package main

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		random, err := local.NewCommand(ctx, "my-bucket", &local.CommandInput{
			Create: pulumi.String("openssl rand -hex 16"),
		})
		if err != nil {
			return err
		}

		ctx.Export("output", random.Stdout)
		return nil
	})
}
```

### Remote provisioning of an EC2 instance

This example creates and EC2 instance, and then uses `remote.Command` and `remote.CopyFile` to run commands and copy files to the remote instance (via SSH). Similar things are possible with Azure, Google Cloud and other cloud provider virtual machines.  Support for Windows-based VMs is being tracked [here](https://github.com/pulumi/pulumi-command/issues/15).

Note that implicit and explicit (`dependsOn`) dependencies can be used to control the order that these `Command` and `CopyFile` resources are constructed relative to each other and to the cloud resources they depend on.  This ensures that the `create` operations run after all dependencies are created, and the `delete` operations run before all dependencies are deleted.

Because the `Command` and `CopyFile` resources replace on changes to their connection, if the EC2 instance is replaced, the commands will all re-run on the new instance (and the `delete` operations will run on the old instance).

Note also that `deleteBeforeReplace` can be composed with `Command` resources to ensure that the `delete` operation on an "old" instance is run before the `create` operation of the new instance, in case a scarce resource is managed by the command.  Similarly, other resource options can naturally be applied to `Command` resources, like `ignoreChanges`.

```ts
import { interpolate, Config } from "@pulumi/pulumi";
import { local, remote, types } from "@pulumi/command";
import * as aws from "@pulumi/aws";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import { size } from "./size";

const config = new Config();
const keyName = config.get("keyName") ?? new aws.ec2.KeyPair("key", { publicKey: config.require("publicKey") }).keyName;
const privateKeyBase64 = config.get("privateKeyBase64");
const privateKey = privateKeyBase64 ? Buffer.from(privateKeyBase64, 'base64').toString('ascii') : fs.readFileSync(path.join(os.homedir(), ".ssh", "id_rsa")).toString("utf8");

const secgrp = new aws.ec2.SecurityGroup("secgrp", {
    description: "Foo",
    ingress: [
        { protocol: "tcp", fromPort: 22, toPort: 22, cidrBlocks: ["0.0.0.0/0"] },
        { protocol: "tcp", fromPort: 80, toPort: 80, cidrBlocks: ["0.0.0.0/0"] },
    ],
});

const ami = aws.ec2.getAmiOutput({
    owners: ["amazon"],
    mostRecent: true,
    filters: [{
        name: "name",
        values: ["amzn2-ami-hvm-2.0.????????-x86_64-gp2"],
    }],
});

const server = new aws.ec2.Instance("server", {
    instanceType: size,
    ami: ami.id,
    keyName: keyName,
    vpcSecurityGroupIds: [secgrp.id],
}, { replaceOnChanges: ["instanceType"] });

// Now set up a connection to the instance and run some provisioning operations on the instance.

const connection: types.input.remote.ConnectionInput = {
    host: server.publicIp,
    user: "ec2-user",
    privateKey: privateKey,
};

const hostname = new remote.Command("hostname", {
    connection,
    create: "hostname",
});

new remote.Command("remotePrivateIP", {
    connection,
    create: interpolate`echo ${server.privateIp} > private_ip.txt`,
    delete: `rm private_ip.txt`,
}, { deleteBeforeReplace: true });

new local.Command("localPrivateIP", {
    create: interpolate`echo ${server.privateIp} > private_ip.txt`,
    delete: `rm private_ip.txt`,
}, { deleteBeforeReplace: true });

const sizeFile = new remote.CopyFile("size", {
    connection,
    localPath: "./size.ts",
    remotePath: "size.ts",
})

const catSize = new remote.Command("checkSize", {
    connection,
    create: "cat size.ts",
}, { dependsOn: sizeFile })

export const confirmSize = catSize.stdout;
export const publicIp = server.publicIp;
export const publicHostName = server.publicDns;
export const hostnameStdout = hostname.stdout;
```

### Invoking a Lambda during Pulumi deployment

There may be cases where it is useful to run some code within an AWS Lambda or other serverless function during the deployment.  For example, this may allow running some code from within a VPC, or with a specific role, without needing to have persistent compute available (such as the EC2 example above).

Note that the Lambda function itself can be created within the same Pulumi program, and then invoked after creation.  

The example below simply creates some random value within the Lambda, which is a very roundabout way of doing the same thing as the first "random" example above, but this pattern can be used for more complex scenarios where the Lambda does things a local script could not.

```ts
import { local } from "@pulumi/command";
import * as aws from "@pulumi/aws";
import * as crypto from "crypto";

const f = new aws.lambda.CallbackFunction("f", {
    publish: true,
    callback: async (ev: any) => {
        return crypto.randomBytes(ev.len/2).toString('hex');
    }
});

const rand = new local.Command("execf", {
    create: `aws lambda invoke --function-name "$FN" --payload '{"len": 10}' --cli-binary-format raw-in-base64-out out.txt >/dev/null && cat out.txt | tr -d '"'  && rm out.txt`,
    environment: {
        FN: f.qualifiedArn,
        AWS_REGION: aws.config.region!,
        AWS_PAGER: "",
    },
})

export const output = rand.stdout;
```

### Using `local.Command `with CURL to manage external REST API

This example uses `local.Command` to create a simple resource provider for managing GitHub labels, by invoking `curl` commands on `create` and `delete` commands against the GitHub REST API.  A similar approach could be applied to build other simple providers against any REST API directly from within Pulumi programs in any language.  This approach is somewhat limited by the fact that `local.Command` does not yet support `diff`/`update`/`read`.  Support for those may be [added in the future](https://github.com/pulumi/pulumi-command/issues/20).

This example also shows how `local.Command` can be used as an implementation detail inside a nicer abstraction, like the `GitHubLabel` component defined below.

```ts
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
```

```sh
# create_label.sh
curl \
  -s \
  -X POST \
  -H "authorization: Bearer $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$OWNER/$REPO/labels \
  -d "{\"name\":\"$NAME\"}"
```

```sh
# delete_label.sh
curl \
  -s \
  -X DELETE \
  -H "authorization: Bearer $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$OWNER/$REPO/labels/$NAME
```

### Graceful cleanup of workloads in a Kubernetes cluster

There are cases where it's important to run some cleanup operation before destroying a resource such as when destroying the resource does not properly handle orderly cleanup.  For example, destroying an EKS Cluster will not ensure that all Kubernetes object finalizers are run, which may lead to leaking external resources managed by those Kubernetes resources.  This example shows how we can use a `delete`-only `Command` to ensure some cleanup is run within a cluster before destroying it.

```ts
import { local } from "@pulumi/command";
import * as eks from "@pulumi/eks";
import * as random from "@pulumi/random";
import { interpolate } from "@pulumi/pulumi";

const cluster = new eks.Cluster("cluster", {});

// We could also use `RemoteCommand` to run this from within a node in the cluster
const cleanupKubernetesNamespaces = new local.Command("cleanupKubernetesNamespaces", {
    // This will run before the cluster is destroyed.  Everything else will need to 
    // depend on this resource to ensure this cleanup doesn't happen too early.
    delete: "kubectl delete --all namespaces",
    environment: {
        KUBECONFIG: cluster.kubeconfig,
    },
});
```

### Working with Assets and Paths

When a local command creates assets as part of its execution, these can be captured by specifying `assetPaths` or `archivePaths`.

```typescript
const lambdaBuild = local.runOutput({
    dir: "../my-function",
    command: `yarn && yarn build`,
    archivePaths: ["dist/**"],
});

new aws.lambda.Function("my-function", {
    code: lambdaBuild.archive,
   // ...
});
```

When using the `assetPaths` and `archivePaths`, they take a list of 'globs'.
- We only include files not directories for assets and archives.
- Path separators are `/` on all platforms - including Windows.
- Patterns starting with `!` are 'exclude' rules.
- Rules are evaluated in order, so exclude rules should be after inclusion rules.
- `*` matches anything except `/`
- `**` matches anything, _including_ `/`
- All returned paths are relative to the working directory (without leading `./`) e.g. `file.text` or `subfolder/file.txt`.
- For full details of the globbing syntax, see [github.com/gobwas/glob](https://github.com/gobwas/glob)

#### Asset Paths Example

Given the rules:
```yaml
- "assets/**"
- "src/**.js"
- "!**secret.*"
```

When evaluating against this folder:

```yaml
- assets/
  - logos/
    - logo.svg
- src/
  - index.js
  - secret.js
```

The following paths will be returned:

```yaml
- assets/logos/logo.svg
- src/index.js
```

## Building

### Dependencies

- Go 1.17
- NodeJS 10.X.X or later
- Python 3.6 or later
- .NET Core 3.1

Please refer to [Contributing to Pulumi](https://github.com/pulumi/pulumi/blob/master/CONTRIBUTING.md) for installation
guidance.

### Building locally

Run the following commands to install Go modules, generate all SDKs, and build the provider:

```
$ make ensure
$ make build
$ make install
```

Add the `bin` folder to your `$PATH` or copy the `bin/pulumi-resource-aws-native` file to another location in your `$PATH`.

### Running an example

Navigate to the simple example and run Pulumi:

```
$ cd examples/simple
$ yarn link @pulumi/command
$ yarn install
$ pulumi up
```

