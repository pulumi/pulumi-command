// Copyright 2021, Pulumi Corporation.  All rights reserved.

package examples

import (
	"fmt"
	"os"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

func getRegion(t *testing.T) string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-2"
		fmt.Println("Defaulting region to 'us-east-2'. You can override using the AWS_REGION variable.")
	}

	return region
}

func getBaseOptions(t *testing.T) integration.ProgramTestOptions {
	awsRegion := getRegion(t)
	return integration.ProgramTestOptions{
		ExpectRefreshChanges: true,
		Config: map[string]string{
			"aws:region": awsRegion,
		},
	}
}

func getCwd(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.FailNow()
	}

	return cwd
}

func skipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
}
