package main

import (
	"fmt"
	"math/rand"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
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
		connection := remote.ConnectionArgs{
			Host:     pulumi.String("192.168.3.201"),
			User:     pulumi.StringPtr("root"),
			Password: pulumi.StringPtr("un2Trois$"),
		}
		random, err := remote.NewCommand(ctx, "my-bucket-z", &remote.CommandArgs{
			Connection:             connection,
			Create:                 pulumi.String(fmt.Sprintf("RAND=%s openssl rand -hex 100000 | tr 'a' '\n'", RandStringBytes(32))),
			AddPreviousOutputInEnv: pulumi.BoolPtr(false),
		})
		if err != nil {
			return err
		}

		ctx.Export("output", random.Stdout)
		return nil
	})
}
