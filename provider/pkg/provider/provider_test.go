// Copyright 2024, Pulumi Corporation.  All rights reserved.

package provider

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pulumi/providertest"
	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/assertpreview"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

const baselineVersion = "0.11.1"

func TestUpgradeLocalCommand(t *testing.T) {
	t.Run("stdin", func(t *testing.T) {
		test(t, "stdin")
	})
}

func test(t *testing.T, example string) {
	t.Parallel()

	dir := fmt.Sprintf("../../../examples/%s", example)

	test := pulumitest.NewPulumiTest(t, dir,
		opttest.YarnLink("@pulumi/command"),
		opttest.LocalProviderPath("command", filepath.Join(dir, "../..", "bin")),
	)
	result := providertest.PreviewProviderUpgrade(t, test, "command", baselineVersion)
	assertpreview.HasNoChanges(t, result)
}
