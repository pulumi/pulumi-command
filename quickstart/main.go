package main

import (
	"fmt"
	"math/rand"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		random, err := local.NewCommand(ctx, "my-bucket", &local.CommandArgs{
			Create: pulumi.String(fmt.Sprintf("RAND=%s openssl rand -hex 100000 | tr 'a' '\n'", RandStringBytes(32))),
		})
		if err != nil {
			return err
		}

		ctx.Export("output", random.Stdout)
		return nil
	})
}
