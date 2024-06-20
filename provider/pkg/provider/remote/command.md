A command to run on a remote host. The connection is established via ssh.

{{% examples %}}
## Example Usage

{{% example %}}
### Triggers

This example defines several trigger values of various kinds. Changes to any of them will cause `cmd` to be re-run.

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as command from "@pulumi/command";
import * as random from "@pulumi/random";

const str = "foo";
const fileAsset = new pulumi.asset.FileAsset("Pulumi.yaml");
const rand = new random.RandomString("rand", {length: 5});
const localFile = new command.local.Command("localFile", {
    create: "touch foo.txt",
    archivePaths: ["*.txt"],
});
const cmd = new command.remote.Command("cmd", {
    connection: {
        host: "insert host here",
    },
    create: "echo create > op.txt",
    delete: "echo delete >> op.txt",
    triggers: [
        str,
        rand.result,
        fileAsset,
        localFile.archive,
    ],
});

```

```python
import pulumi
import pulumi_command as command
import pulumi_random as random

foo = "foo"
file_asset_var = pulumi.FileAsset("Pulumi.yaml")
rand = random.RandomString("rand", length=5)
local_file = command.local.Command("localFile",
    create="touch foo.txt",
    archive_paths=["*.txt"])

cmd = command.remote.Command("cmd",
    connection=command.remote.ConnectionArgs(
        host="insert host here",
    ),
    create="echo create > op.txt",
    delete="echo delete >> op.txt",
    triggers=[
        foo,
        rand.result,
        file_asset_var,
        local_file.archive,
    ])
```

```go
package main

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		str := pulumi.String("foo")

		fileAsset := pulumi.NewFileAsset("Pulumi.yaml")

		rand, err := random.NewRandomString(ctx, "rand", &random.RandomStringArgs{
			Length: pulumi.Int(5),
		})
		if err != nil {
			return err
		}

		localFile, err := local.NewCommand(ctx, "localFile", &local.CommandArgs{
			Create: pulumi.String("touch foo.txt"),
			ArchivePaths: pulumi.StringArray{
				pulumi.String("*.txt"),
			},
		})
		if err != nil {
			return err
		}

		_, err = remote.NewCommand(ctx, "cmd", &remote.CommandArgs{
			Connection: &remote.ConnectionArgs{
				Host: pulumi.String("insert host here"),
			},
			Create: pulumi.String("echo create > op.txt"),
			Delete: pulumi.String("echo delete >> op.txt"),
			Triggers: pulumi.Array{
				str,
				rand.Result,
				fileAsset,
				localFile.Archive,
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```

```csharp
using Pulumi;
using Command = Pulumi.Command;
using Random = Pulumi.Random;

return await Deployment.RunAsync(() => 
{
    var str = "foo";

    var fileAssetVar = new FileAsset("Pulumi.yaml");

    var rand = new Random.RandomString("rand", new()
    {
        Length = 5,
    });

    var localFile = new Command.Local.Command("localFile", new()
    {
        Create = "touch foo.txt",
        ArchivePaths = new[]
        {
            "*.txt",
        },
    });

    var cmd = new Command.Remote.Command("cmd", new()
    {
        Connection = new Command.Remote.Inputs.ConnectionArgs
        {
            Host = "insert host here",
        },
        Create = "echo create > op.txt",
        Delete = "echo delete >> op.txt",
        Triggers = new object[]
        {
            str,
            rand.Result,
            fileAssetVar,
            localFile.Archive,
        },
    });

});
```

```java
public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        final var fileAssetVar = new FileAsset("Pulumi.yaml");

        var rand = new RandomString("rand", RandomStringArgs.builder()
            .length(5)
            .build());

        var localFile = new Command("localFile", CommandArgs.builder()
            .create("touch foo.txt")
            .archivePaths("*.txt")
            .build());

        var cmd = new Command("cmd", CommandArgs.builder()
            .connection(ConnectionArgs.builder()
                .host("insert host here")
                .build())
            .create("echo create > op.txt")
            .delete("echo delete >> op.txt")
            .triggers(            
                rand.result(),
                fileAssetVar,
                localFile.archive())
            .build());

    }
}
```

```yaml
config: {}
outputs: {}

resources:
  rand:
    type: random:index/randomString:RandomString
    properties:
      length: 5

  localFile:
    type: command:local:Command
    properties:
      create: touch foo.txt
      archivePaths:
        - "*.txt"

  cmd:
    type: command:remote:Command
    properties:
      connection:
        host: "insert host here"
      create: echo create > op.txt
      delete: echo delete >> op.txt
      triggers:
        - ${rand.result}
        - ${fileAsset}
        - ${localFile.archive}

variables:
  fileAsset:
    fn::fileAsset: "Pulumi.yaml"
```

{{% /example %}}

{{% /examples %}}