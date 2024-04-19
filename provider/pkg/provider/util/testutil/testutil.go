package testutil

import (
	"bytes"
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

// TestContext is an implementation of p.Context that records all log messages in a buffer, regardless of severity.
type TestContext struct {
	context.Context
	Output bytes.Buffer
}

func (c *TestContext) log(msg string) {
	c.Output.WriteString(msg)
}

func (c *TestContext) Log(_ diag.Severity, msg string)                  { c.log(msg) }
func (c *TestContext) Logf(_ diag.Severity, msg string, _ ...any)       { c.log(msg) }
func (c *TestContext) LogStatus(_ diag.Severity, msg string)            { c.log(msg) }
func (c *TestContext) LogStatusf(_ diag.Severity, msg string, _ ...any) { c.log(msg) }
func (c *TestContext) RuntimeInformation() p.RunInfo                    { return p.RunInfo{} }
