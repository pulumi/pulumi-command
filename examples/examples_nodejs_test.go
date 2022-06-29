// Copyright 2016-2021, Pulumi Corporation.  All rights reserved.
//go:build nodejs || all
// +build nodejs all

package examples

import (
	"encoding/base64"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "random"),
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				out, ok := stack.Outputs["output"].(string)
				assert.True(t, ok)
				assert.Len(t, out, 32)
			},
		})
	integration.ProgramTest(t, &test)
}

func TestStdin(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "stdin"),
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				out, ok := stack.Outputs["output"].(string)
				assert.True(t, ok)
				assert.Equal(t, "the quick brown fox", out)
			},
		})
	integration.ProgramTest(t, &test)
}

func TestSimple(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:                    filepath.Join(getCwd(t), "simple"),
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {},
			EditDirs: []integration.EditDir{
				{
					Dir:      filepath.Join("simple", "update"),
					Additive: true,
					ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
						replaces := 0
						for _, ev := range stack.Events {
							if ev.ResourcePreEvent != nil {
								if ev.ResourcePreEvent.Metadata.Op == apitype.OpReplace {
									replaces++
								}
							}
						}
						assert.Equal(t, 0, replaces)
					},
				},
				{
					Dir:      filepath.Join("simple", "replace"),
					Additive: true,
					ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
						replaces := 0
						for _, ev := range stack.Events {
							if ev.ResourcePreEvent != nil {
								if ev.ResourcePreEvent.Metadata.Op == apitype.OpReplace {
									replaces++
								}
							}
						}
						assert.Equal(t, 4, replaces)
					},
				},
				{
					Dir:           filepath.Join("simple", "fail"),
					Additive:      true,
					ExpectFailure: true,
				},
				{
					Dir:           filepath.Join("simple", "update-change"),
					Additive:      true,
					ExpectFailure: true,
				},
			},
		})
	integration.ProgramTest(t, &test)
}

func TestSimpleWithUpdate(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:                    filepath.Join(getCwd(t), "simple-with-update"),
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {},
			EditDirs: []integration.EditDir{
				{
					Dir:      filepath.Join("simple-with-update", "update-change"),
					Additive: true,
				},
			},
		})
	integration.ProgramTest(t, &test)
}

func TestEc2RemoteTs(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(getRegion(t))},
	)
	assert.NoError(t, err)
	svc := ec2.New(sess)
	keyName, err := resource.NewUniqueHex("test-keyname", 8, 20)
	assert.NoError(t, err)
	t.Logf("Creating keypair %s.\n", keyName)
	key, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(keyName),
	})
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer func() {
		t.Logf("Deleting keypair %s.\n", keyName)
		_, err := svc.DeleteKeyPair(&ec2.DeleteKeyPairInput{
			KeyName: aws.String(keyName),
		})
		assert.NoError(t, err)
	}()
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "ec2_remote"),
			Config: map[string]string{
				"keyName": aws.StringValue(key.KeyName),
			},
			Secrets: map[string]string{
				"privateKeyBase64": base64.StdEncoding.EncodeToString([]byte(aws.StringValue(key.KeyMaterial))),
			},
			EditDirs: []integration.EditDir{{
				Dir:      filepath.Join("ec2_remote", "replace_instance"),
				Additive: true,
			}},
		})

	integration.ProgramTest(t, &test)
}

func TestLambda(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "lambda"),
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				out, ok := stack.Outputs["output"].(string)
				assert.True(t, ok)
				assert.Len(t, out, 10)
			},
		})
	integration.ProgramTest(t, &test)
}

func TestSimpleRun(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "simple-run"),
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				assets, ok := stack.Outputs["plainAssets"].(map[string]interface{})
				assert.True(t, ok)
				assert.Len(t, assets, 1)
			},
		})
	integration.ProgramTest(t, &test)
}

func getJSBaseOptions(t *testing.T) integration.ProgramTestOptions {
	base := getBaseOptions(t)
	baseJS := base.With(integration.ProgramTestOptions{
		Dependencies: []string{
			"@pulumi/command",
		},
	})

	return baseJS
}
