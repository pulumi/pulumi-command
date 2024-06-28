// Copyright 2024, Pulumi Corporation.  All rights reserved.

package provider

import (
	"fmt"
	"testing"

	"github.com/pulumi/providertest"
)

const baselineVersion = "0.11.1"

func TestUpgradeLocalCommand(t *testing.T) {
	t.Run("stdin", func(t *testing.T) {
		runExampleParallel(t, "stdin")
	})
}

func runExampleParallel(t *testing.T, example string, opts ...providertest.Option) {
	t.Parallel()
	test(fmt.Sprintf("../../../examples/%s", example), opts...).Run(t)
}

func test(dir string, opts ...providertest.Option) *providertest.ProviderTest {
	opts = append(opts,
		providertest.WithProviderName("command"),
		providertest.WithBaselineVersion(baselineVersion),
	)
	return providertest.NewProviderTest(dir, opts...)
}
