package util

import (
	"bufio"
	"io"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

func CopyOutput(ctx p.Context, r io.Reader, doneCh chan<- struct{}, severity diag.Severity) {
	defer close(doneCh)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ctx.Log(severity, scanner.Text())
	}
}
