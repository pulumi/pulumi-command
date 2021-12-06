// Copyright 2016-2021, Pulumi Corporation.  All rights reserved.
//go:build nodejs || all
// +build nodejs all

package examples

import (
	"encoding/base64"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestSimpleTs(t *testing.T) {
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

func getJSBaseOptions(t *testing.T) integration.ProgramTestOptions {
	base := getBaseOptions(t)
	baseJS := base.With(integration.ProgramTestOptions{
		Dependencies: []string{
			"@pulumi/command",
		},
	})

	return baseJS
}
