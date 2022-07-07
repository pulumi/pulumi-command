// Copyright 2016-2021, Pulumi Corporation.  All rights reserved.
//go:build go || all
// +build go all

package examples

import (
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestRandomGo(t *testing.T) {
	test := getGoBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "random-go"),
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				out, ok := stack.Outputs["output"].(string)
				assert.True(t, ok)
				assert.Len(t, out, 32)
			},
		})
	integration.ProgramTest(t, &test)
}

func TestStdinGo(t *testing.T) {
	test := getGoBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "stdin-go"),
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				out, ok := stack.Outputs["output"].(string)
				assert.True(t, ok)
				assert.Equal(t, "the quick brown fox", out)
				assets, ok := stack.Outputs["assets"].(map[string]interface{})
				assert.True(t, ok)
				assert.Len(t, assets, 1)
			},
		})
	integration.ProgramTest(t, &test)
}

func getGoBaseOptions(t *testing.T) integration.ProgramTestOptions {
	base := getBaseOptions(t)
	baseGo := base.With(integration.ProgramTestOptions{
		Verbose: true,
		Dependencies: []string{
			"github.com/pulumi/pulumi-command/sdk",
		},
	})

	return baseGo
}
