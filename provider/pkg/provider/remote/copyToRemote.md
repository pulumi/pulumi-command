Copy an Asset or Archive to a remote host.

{{% examples %}}

## Example usage

This example copies a local directory to a remote host via SSH. For brevity, the remote server is assumed to exist, but it could also be provisioned in the same Pulumi program.

{{% example %}}

```typescript
import * as pulumi from "@pulumi/pulumi";
import { remote, types } from "@pulumi/command";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

export = async () => {
    const config = new pulumi.Config();

    // Get the private key to connect to the server. If a key is
    // provided, use it, otherwise default to the standard id_rsa SSH key.
    const privateKeyBase64 = config.get("privateKeyBase64");
    const privateKey = privateKeyBase64 ?
        Buffer.from(privateKeyBase64, 'base64').toString('ascii') :
        fs.readFileSync(path.join(os.homedir(), ".ssh", "id_rsa")).toString("utf8");

    const serverPublicIp = config.require("serverPublicIp");
    const userName = config.require("userName");

    // The configuration of our SSH connection to the instance.
    const connection: types.input.remote.ConnectionArgs = {
        host: serverPublicIp,
        user: userName,
        privateKey: privateKey,
    };

    // Set up source and target of the remote copy.
    const from = config.require("payload")!;
    const archive = new pulumi.asset.FileArchive(from);
    const to = config.require("destDir")!;

    // Copy the files to the remote.
    const copy = new remote.CopyToRemote("copy", {
        connection,
        source: archive,
        remotePath: to,
    });

    // Verify that the expected files were copied to the remote.
    // We want to run this after each copy, i.e., when something changed,
    // so we use the asset to be copied as a trigger.
    const find = new remote.Command("ls", {
        connection,
        create: `find ${to}/${from} | sort`,
        triggers: [archive],
    }, { dependsOn: copy });

    return {
        remoteContents: find.stdout
    }
}
```

```python
import pulumi
import pulumi_command as command

config = pulumi.Config()

server_public_ip = config.require("serverPublicIp")
user_name = config.require("userName")
private_key = config.require("privateKey")
payload = config.require("payload")
dest_dir = config.require("destDir")

archive = pulumi.FileArchive(payload)

# The configuration of our SSH connection to the instance.
conn = command.remote.ConnectionArgs(
    host = server_public_ip,
    user = user_name,
    privateKey = private_key,
)

# Copy the files to the remote.
copy = command.remote.CopyToRemote("copy",
    connection=conn,
    source=archive,
    remote_path=dest_dir)

# Verify that the expected files were copied to the remote.
# We want to run this after each copy, i.e., when something changed,
# so we use the asset to be copied as a trigger.
find = command.remote.Command("find",
    connection=conn,
    create=f"find {dest_dir}/{payload} | sort",
    triggers=[archive],
    opts = pulumi.ResourceOptions(depends_on=[copy]))

pulumi.export("remoteContents", find.stdout)
```

```go
package main

import (
	"fmt"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		serverPublicIp := cfg.Require("serverPublicIp")
		userName := cfg.Require("userName")
		privateKey := cfg.Require("privateKey")
		payload := cfg.Require("payload")
		destDir := cfg.Require("destDir")

		archive := pulumi.NewFileArchive(payload)

		conn := remote.ConnectionArgs{
			Host:       pulumi.String(serverPublicIp),
			User:       pulumi.String(userName),
			PrivateKey: pulumi.String(privateKey),
		}

		copy, err := remote.NewCopyToRemote(ctx, "copy", &remote.CopyToRemoteArgs{
			Connection: conn,
			Source:     archive,
		})
		if err != nil {
			return err
		}

		find, err := remote.NewCommand(ctx, "find", &remote.CommandArgs{
			Connection: conn,
			Create:     pulumi.String(fmt.Sprintf("find %v/%v | sort", destDir, payload)),
			Triggers: pulumi.Array{
				archive,
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			copy,
		}))
		if err != nil {
			return err
		}

		ctx.Export("remoteContents", find.Stdout)
		return nil
	})
}
```

```csharp
using System.Collections.Generic;
using Pulumi;
using Command = Pulumi.Command;

return await Deployment.RunAsync(() =>
{
    var config = new Config();
    var serverPublicIp = config.Require("serverPublicIp");
    var userName = config.Require("userName");
    var privateKey = config.Require("privateKey");
    var payload = config.Require("payload");
    var destDir = config.Require("destDir");

    var archive = new FileArchive(payload);

    // The configuration of our SSH connection to the instance.
    var conn = new Command.Remote.Inputs.ConnectionArgs
    {
        Host = serverPublicIp,
        User = userName,
        PrivateKey = privateKey,
    };

    // Copy the files to the remote.
    var copy = new Command.Remote.CopyToRemote("copy", new()
    {
        Connection = conn,
        Source = archive,
    });

    // Verify that the expected files were copied to the remote.
    // We want to run this after each copy, i.e., when something changed,
    // so we use the asset to be copied as a trigger.
    var find = new Command.Remote.Command("find", new()
    {
        Connection = conn,
        Create = $"find {destDir}/{payload} | sort",
        Triggers = new[]
        {
            archive,
        },
    }, new CustomResourceOptions
    {
        DependsOn =
        {
            copy,
        },
    });

    return new Dictionary<string, object?>
    {
        ["remoteContents"] = find.Stdout,
    };
});
```

```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.command.remote.Command;
import com.pulumi.command.remote.CommandArgs;
import com.pulumi.command.remote.CopyToRemote;
import com.pulumi.command.remote.inputs.*;
import com.pulumi.resources.CustomResourceOptions;
import com.pulumi.asset.FileArchive;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        final var config = ctx.config();
        final var serverPublicIp = config.require("serverPublicIp");
        final var userName = config.require("userName");
        final var privateKey = config.require("privateKey");
        final var payload = config.require("payload");
        final var destDir = config.require("destDir");

        final var archive = new FileArchive(payload);

        // The configuration of our SSH connection to the instance.
        final var conn = ConnectionArgs.builder()
            .host(serverPublicIp)
            .user(userName)
            .privateKey(privateKey)
            .build();

        // Copy the files to the remote.
        var copy = new CopyToRemote("copy", CopyToRemoteArgs.builder()
            .connection(conn)
            .source(archive)
            .destination(destDir)
            .build());

        // Verify that the expected files were copied to the remote.
        // We want to run this after each copy, i.e., when something changed,
        // so we use the asset to be copied as a trigger.
        var find = new Command("find", CommandArgs.builder()
            .connection(conn)
            .create(String.format("find %s/%s | sort", destDir,payload))
            .triggers(archive)
            .build(), CustomResourceOptions.builder()
                .dependsOn(copy)
                .build());

        ctx.export("remoteContents", find.stdout());
    }
}
```

```yaml
resources:
  # Copy the files to the remote.
  copy:
    type: command:remote:CopyToRemote
    properties:
      connection: ${conn}
      source: ${archive}
      remotePath: ${destDir}

  # Verify that the expected files were copied to the remote.
  # We want to run this after each copy, i.e., when something changed,
  # so we use the asset to be copied as a trigger.
  find:
    type: command:remote:Command
    properties:
      connection: ${conn}
      create: find ${destDir}/${payload} | sort
      triggers:
        - ${archive}
    options:
      dependsOn:
        - ${copy}

config:
  serverPublicIp:
    type: string
  userName:
    type: string
  privateKey:
    type: string
  payload:
    type: string
  destDir:
    type: string

variables:
  # The source directory or archive to copy.
  archive:
    fn::fileArchive: ${payload}
  # The configuration of our SSH connection to the instance.
  conn:
    host: ${serverPublicIp}
    user: ${userName}
    privateKey: ${privateKey}

outputs:
  remoteContents: ${find.stdout}
```

{{% /example %}}

{{% /examples %}}

