A local command to be executed.

This command can be inserted into the life cycles of other resources using the `dependsOn` or `parent` resource options. A command is considered to have failed when it finished with a non-zero exit code. This will fail the CRUD step of the `Command` resource.

{{% examples %}}
## Example Usage

{{% example %}}
### Triggers

This example defines several trigger values of various kinds. Changes to any of them will cause `cmd` to be re-run. However, note that for `fileAsset` it's the variable itself that is the trigger, not the contents of index.ts, since triggers are simply opaque values.

```typescript
import * as local from "@pulumi/command/local";
import * as random from "@pulumi/random";
import { asset } from "@pulumi/pulumi";
import * as path from "path";

const str = "foo";
const fileAsset = new pulumi.asset.FileAsset("Pulumi.yaml");
const rand = new random.RandomString("rand", {length: 5});
const localFile = new command.local.Command("localFile", {
    create: "touch foo.txt",
    archivePaths: ["*.txt"],
});

const cmd = new local.Command("pwd", {
    create: "echo create > op.txt",
    delete: "echo delete >> op.txt",
    triggers: [str, rand.result, fileAsset, localFile.archive],
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

cmd = command.local.Command("cmd",
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

		_, err = local.NewCommand(ctx, "cmd", &local.CommandArgs{
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

    var cmd = new Command.Local.Command("cmd", new()
    {
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
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.random.RandomString;
import com.pulumi.random.RandomStringArgs;
import com.pulumi.command.local.Command;
import com.pulumi.command.local.CommandArgs;
import com.pulumi.asset.FileAsset;

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
            .create("echo create > op.txt")
            .delete("echo delete >> op.txt")
            .triggers(            
                "foo",
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
    type: command:local:Command
    properties:
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