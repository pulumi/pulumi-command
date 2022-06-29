package main

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		random, err := local.NewCommand(ctx, "stdin", &local.CommandArgs{
			Create:     pulumi.String("head -n 1"),
			Stdin:      pulumi.String("the quick brown fox\njumped over\nthe lazy dog"),
			AssetPaths: pulumi.StringArray{pulumi.String("*.go")},
		})
		if err != nil {
			return err
		}

		ctx.Export("output", random.Stdout)
		ctx.Export("assets", random.Assets)
		return nil
	})
}
