A local command to be executed.

This command can be inserted into the life cycles of other resources using the `dependsOn` or `parent` resource options. A command is considered to have failed when it finished with a non-zero exit code. This will fail the CRUD step of the `Command` resource.

{{% examples %}}

## Example Usage

{{% example %}}

### Basic Example

This example shows the simplest use case, simply running a command on `create` in the Pulumi lifecycle.

```typescript
import { local } from "@pulumi/command";

const random = new local.Command("random", {
    create: "openssl rand -hex 16",
});

export const output = random.stdout;
```

```csharp
using System.Collections.Generic;
using Pulumi;
using Pulumi.Command.Local;

await Deployment.RunAsync(() =>
{
    var command = new Command("random", new CommandArgs
    {
        Create = "openssl rand -hex 16"
    });

    return new Dictionary<string, object?>
    {
        ["stdOut"] = command.Stdout
    };
});
```

```python
import pulumi
from pulumi_command import local

random = local.Command("random",
    create="openssl rand -hex 16"
)

pulumi.export("random", random.stdout)
```

```go
package main

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		random, err := local.NewCommand(ctx, "my-bucket", &local.CommandArgs{
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

```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.command.local.Command;
import com.pulumi.command.local.CommandArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        var random = new Command("random", CommandArgs.builder()
            .create("openssl rand -hex 16")
            .build());

        ctx.export("rand", random.stdout());
    }
}
```

```yaml
outputs:
  rand: "${random.stdout}"
resources:
  random:
    type: command:local:Command
    properties:
      create: "openssl rand -hex 16"
```

{{% /example %}}

{{% example %}}

### Invoking a Lambda during Pulumi Deployment

This example show using a local command to invoke an AWS Lambda once it's deployed. The Lambda invocation could also depend on other resources.

```typescript
import * as aws from "@pulumi/aws";
import { local } from "@pulumi/command";
import { getStack } from "@pulumi/pulumi";

const f = new aws.lambda.CallbackFunction("f", {
    publish: true,
    callback: async (ev: any) => {
        return `Stack ${ev.stackName} is deployed!`;
    }
});

const invoke = new local.Command("execf", {
    create: `aws lambda invoke --function-name "$FN" --payload '{"stackName": "${getStack()}"}' --cli-binary-format raw-in-base64-out out.txt >/dev/null && cat out.txt | tr -d '"'  && rm out.txt`,
    environment: {
        FN: f.qualifiedArn,
        AWS_REGION: aws.config.region!,
        AWS_PAGER: "",
    },
}, { dependsOn: f })

export const output = invoke.stdout;
```

```python
import pulumi
import json
import pulumi_aws as aws
import pulumi_command as command

lambda_role = aws.iam.Role("lambdaRole", assume_role_policy=json.dumps({
    "Version": "2012-10-17",
    "Statement": [{
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
            "Service": "lambda.amazonaws.com",
        },
    }],
}))

lambda_function = aws.lambda_.Function("lambdaFunction",
    name="f",
    publish=True,
    role=lambda_role.arn,
    handler="index.handler",
    runtime=aws.lambda_.Runtime.NODE_JS20D_X,
    code=pulumi.FileArchive("./handler"))

aws_config = pulumi.Config("aws")
aws_region = aws_config.require("region")

invoke_command = command.local.Command("invokeCommand",
    create=f"aws lambda invoke --function-name \"$FN\" --payload '{{\"stackName\": \"{pulumi.get_stack()}\"}}' --cli-binary-format raw-in-base64-out out.txt >/dev/null && cat out.txt | tr -d '\"'  && rm out.txt",
    environment={
        "FN": lambda_function.arn,
        "AWS_REGION": aws_region,
        "AWS_PAGER": "",
    },
    opts = pulumi.ResourceOptions(depends_on=[lambda_function]))

pulumi.export("output", invoke_command.stdout)
```

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		awsConfig := config.New(ctx, "aws")
		awsRegion := awsConfig.Require("region")

		tmpJSON0, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": "lambda.amazonaws.com",
					},
				},
			},
		})
		if err != nil {
			return err
		}
		json0 := string(tmpJSON0)
		lambdaRole, err := iam.NewRole(ctx, "lambdaRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(json0),
		})
		if err != nil {
			return err
		}
		lambdaFunction, err := lambda.NewFunction(ctx, "lambdaFunction", &lambda.FunctionArgs{
			Name:    pulumi.String("f"),
			Publish: pulumi.Bool(true),
			Role:    lambdaRole.Arn,
			Handler: pulumi.String("index.handler"),
			Runtime: pulumi.String(lambda.RuntimeNodeJS20dX),
			Code:    pulumi.NewFileArchive("./handler"),
		})
		if err != nil {
			return err
		}
		invokeCommand, err := local.NewCommand(ctx, "invokeCommand", &local.CommandArgs{
			Create: pulumi.String(fmt.Sprintf("aws lambda invoke --function-name \"$FN\" --payload '{\"stackName\": \"%v\"}' --cli-binary-format raw-in-base64-out out.txt >/dev/null && cat out.txt | tr -d '\"'  && rm out.txt", ctx.Stack())),
			Environment: pulumi.StringMap{
				"FN":         lambdaFunction.Arn,
				"AWS_REGION": pulumi.String(awsRegion),
				"AWS_PAGER":  pulumi.String(""),
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			lambdaFunction,
		}))
		if err != nil {
			return err
		}
		ctx.Export("output", invokeCommand.Stdout)
		return nil
	})
}
```

```csharp
using System.Collections.Generic;
using System.Text.Json;
using Pulumi;
using Aws = Pulumi.Aws;
using Command = Pulumi.Command;

return await Deployment.RunAsync(() => 
{
    var awsConfig = new Config("aws");

    var lambdaRole = new Aws.Iam.Role("lambdaRole", new()
    {
        AssumeRolePolicy = JsonSerializer.Serialize(new Dictionary<string, object?>
        {
            ["Version"] = "2012-10-17",
            ["Statement"] = new[]
            {
                new Dictionary<string, object?>
                {
                    ["Action"] = "sts:AssumeRole",
                    ["Effect"] = "Allow",
                    ["Principal"] = new Dictionary<string, object?>
                    {
                        ["Service"] = "lambda.amazonaws.com",
                    },
                },
            },
        }),
    });

    var lambdaFunction = new Aws.Lambda.Function("lambdaFunction", new()
    {
        Name = "f",
        Publish = true,
        Role = lambdaRole.Arn,
        Handler = "index.handler",
        Runtime = Aws.Lambda.Runtime.NodeJS20dX,
        Code = new FileArchive("./handler"),
    });

    var invokeCommand = new Command.Local.Command("invokeCommand", new()
    {
        Create = $"aws lambda invoke --function-name \"$FN\" --payload '{{\"stackName\": \"{Deployment.Instance.StackName}\"}}' --cli-binary-format raw-in-base64-out out.txt >/dev/null && cat out.txt | tr -d '\"'  && rm out.txt",
        Environment = 
        {
            { "FN", lambdaFunction.Arn },
            { "AWS_REGION", awsConfig.Require("region") },
            { "AWS_PAGER", "" },
        },
    }, new CustomResourceOptions
    {
        DependsOn =
        {
            lambdaFunction,
        },
    });

    return new Dictionary<string, object?>
    {
        ["output"] = invokeCommand.Stdout,
    };
});
```

```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.aws.iam.Role;
import com.pulumi.aws.iam.RoleArgs;
import com.pulumi.aws.lambda.Function;
import com.pulumi.aws.lambda.FunctionArgs;
import com.pulumi.command.local.Command;
import com.pulumi.command.local.CommandArgs;
import static com.pulumi.codegen.internal.Serialization.*;
import com.pulumi.resources.CustomResourceOptions;
import com.pulumi.asset.FileArchive;
import java.util.Map;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        var awsConfig = ctx.config("aws");
        var awsRegion = awsConfig.require("region");

        var lambdaRole = new Role("lambdaRole", RoleArgs.builder()
                .assumeRolePolicy(serializeJson(
                        jsonObject(
                                jsonProperty("Version", "2012-10-17"),
                                jsonProperty("Statement", jsonArray(jsonObject(
                                        jsonProperty("Action", "sts:AssumeRole"),
                                        jsonProperty("Effect", "Allow"),
                                        jsonProperty("Principal", jsonObject(
                                                jsonProperty("Service", "lambda.amazonaws.com")))))))))
                .build());

        var lambdaFunction = new Function("lambdaFunction", FunctionArgs.builder()
                .name("f")
                .publish(true)
                .role(lambdaRole.arn())
                .handler("index.handler")
                .runtime("nodejs20.x")
                .code(new FileArchive("./handler"))
                .build());

        // Work around the lack of Output.all for Maps in Java. We cannot use a plain Map because
        // `lambdaFunction.arn()` is an Output<String>.
        var invokeEnv = Output.tuple(
                Output.of("FN"), lambdaFunction.arn(),
                Output.of("AWS_REGION"), Output.of(awsRegion),
                Output.of("AWS_PAGER"), Output.of("")
        ).applyValue(t -> Map.of(t.t1, t.t2, t.t3, t.t4, t.t5, t.t6));

        var invokeCommand = new Command("invokeCommand", CommandArgs.builder()
                .create(String.format(
                        "aws lambda invoke --function-name \"$FN\" --payload '{\"stackName\": \"%s\"}' --cli-binary-format raw-in-base64-out out.txt >/dev/null && cat out.txt | tr -d '\"'  && rm out.txt",
                        ctx.stackName()))
                .environment(invokeEnv)
                .build(),
                CustomResourceOptions.builder()
                        .dependsOn(lambdaFunction)
                        .build());

        ctx.export("output", invokeCommand.stdout());
    }
}
```

```yaml
resources:
  lambdaRole:
    type: aws:iam:Role
    properties:
      assumeRolePolicy:
        fn::toJSON:
          Version: "2012-10-17"
          Statement:
            - Action: sts:AssumeRole
              Effect: Allow
              Principal:
                Service: lambda.amazonaws.com

  lambdaFunction:
    type: aws:lambda:Function
    properties:
      name: f
      publish: true
      role: ${lambdaRole.arn}
      handler: index.handler
      runtime: "nodejs20.x"
      code:
        fn::fileArchive: ./handler

  invokeCommand:
    type: command:local:Command
    properties:
      create: 'aws lambda invoke --function-name "$FN" --payload ''{"stackName": "${pulumi.stack}"}'' --cli-binary-format raw-in-base64-out out.txt >/dev/null && cat out.txt | tr -d ''"''  && rm out.txt'
      environment:
        FN: ${lambdaFunction.arn}
        AWS_REGION: ${aws:region}
        AWS_PAGER: ""
    options:
      dependsOn:
        - ${lambdaFunction}

outputs:
  output: ${invokeCommand.stdout}
```

{{% /example %}}

{{% example %}}

### Using Triggers

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

const cmd = new command.local.Command("cmd", {
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