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
