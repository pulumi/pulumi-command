// Copyright 2016-2021, Pulumi Corporation.  All rights reserved.
//go:build nodejs || all
// +build nodejs all

package examples

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestDeleteFromStdout(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "delete-from-stdout"),
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				out, ok := stack.Outputs["output"].(string)
				assert.True(t, ok)
				_, err := os.Stat(out)
				assert.NoError(t, err)
			},
		})
	integration.ProgramTest(t, &test)
}

func TestStderr(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:                    filepath.Join(getCwd(t), "stderr"),
			SkipPreview:            true,
			SkipRefresh:            true,
			SkipEmptyPreviewUpdate: true,
			ExpectFailure:          true,
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				for _, ev := range stack.Events {
					if ev.DiagnosticEvent != nil {
						switch diag.Severity(ev.DiagnosticEvent.Severity) {
						case diag.Info:
							assert.True(t, ev.DiagnosticEvent.Ephemeral)
						case diag.Error:
							if ev.DiagnosticEvent.URN != "" {
								assert.False(t, ev.DiagnosticEvent.Ephemeral)
								assert.Regexp(t, `^exit status \d+: running`, ev.DiagnosticEvent.Message)
							}
						}
					}
				}
			},
		})
	integration.ProgramTest(t, &test)
}

func TestLambdaTs(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "lambda-ts"),
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

func genEC2KeyPair(t *testing.T) (*ec2.CreateKeyPairOutput, func()) {
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
	require.NoError(t, err)

	cleanup := func() {
		t.Logf("Deleting keypair %s.\n", keyName)
		_, err := svc.DeleteKeyPair(&ec2.DeleteKeyPairInput{
			KeyName: aws.String(keyName),
		})
		assert.NoError(t, err)
	}

	return key, cleanup
}

func testEc2Ts(t *testing.T, targetDir string) {
	key, keyCleanup := genEC2KeyPair(t)
	defer keyCleanup()

	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), targetDir),
			Config: map[string]string{
				"keyName": aws.StringValue(key.KeyName),
			},
			Secrets: map[string]string{
				"privateKeyBase64": base64.StdEncoding.EncodeToString([]byte(aws.StringValue(key.KeyMaterial))),
			},
			EditDirs: []integration.EditDir{{
				Dir:      filepath.Join(targetDir, "replace_instance"),
				Additive: true,
			}},
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				isEncrypted := func(v interface{}) bool {
					m, ok := v.(map[string]interface{})
					if !ok {
						return false
					}
					sigKey := m[resource.SigKey]
					if sigKey == nil {
						return false
					}

					v, vOk := sigKey.(string)
					if !vOk {
						return false
					}

					if v != resource.SecretSig {
						return false
					}

					ciphertext := m["ciphertext"]
					if ciphertext == nil {
						return false
					}

					_, cOk := ciphertext.(string)
					return cOk
				}

				assertEncryptedValue := func(m map[string]interface{}, key string) {
					assert.Truef(t, isEncrypted(m[key]), "%s value should be encrypted", key)
				}
				assertEncryptedValue(stack.Outputs, "connectionSecret")
			},
		})

	integration.ProgramTest(t, &test)
}

func TestEc2RemoteTs(t *testing.T) { testEc2Ts(t, "ec2_remote") }

func TestEc2RemoteProxyTs(t *testing.T) { testEc2Ts(t, "ec2_remote_proxy") }

func TestEc2DirCopy(t *testing.T) {
	key, keyCleanup := genEC2KeyPair(t)
	defer keyCleanup()

	const dest = "/tmp/ec2_dir_copy"
	basePath := filepath.Join(getCwd(t), "ec2_dir_copy")

	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: basePath,
			Config: map[string]string{
				"keyName": aws.StringValue(key.KeyName),
				"destDir": dest,
			},
			Secrets: map[string]string{
				"privateKeyBase64": base64.StdEncoding.EncodeToString([]byte(aws.StringValue(key.KeyMaterial))),
			},
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				remoteLSOutput, ok := stack.Outputs["lsRemote"]
				require.True(t, ok)
				remoteLS, ok := remoteLSOutput.(string)
				require.True(t, ok)

				assert.Equal(t,
					dest+"\n"+
						dest+"/file1\n"+
						dest+"/one\n"+
						dest+"/one/file2\n"+
						dest+"/one/two\n"+
						dest+"/one/two/file3",
					remoteLS)
			},
			EditDirs: []integration.EditDir{
				// There's a new file in src/ so a new copy should be made.
				{
					Dir:             filepath.Join(basePath, "step2"),
					Additive:        true,
					ExpectNoChanges: false,
					ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
						remoteLSOutput, ok := stack.Outputs["lsRemote"]
						require.True(t, ok)
						remoteLS, ok := remoteLSOutput.(string)
						require.True(t, ok)

						assert.Equal(t,
							dest+"\n"+
								dest+"/file1\n"+
								dest+"/newfile\n"+ // added
								dest+"/one\n"+
								dest+"/one/file2\n"+
								dest+"/one/two\n"+
								dest+"/one/two/file3",
							remoteLS)
					},
				},
			},
		})

	integration.ProgramTest(t, &test)
}

func TestEc2CopyFile(t *testing.T) {
	key, keyCleanup := genEC2KeyPair(t)
	defer keyCleanup()

	const dest = "/tmp/TestEc2CopyFile"
	basePath := filepath.Join(getCwd(t), "ec2_copyfile")

	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: basePath,
			Config: map[string]string{
				"keyName": aws.StringValue(key.KeyName),
				"destDir": dest,
			},
			Secrets: map[string]string{
				"privateKeyBase64": base64.StdEncoding.EncodeToString([]byte(aws.StringValue(key.KeyMaterial))),
			},
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				remoteLSOutput, ok := stack.Outputs["lsRemote"]
				require.True(t, ok)
				remoteLS, ok := remoteLSOutput.(string)
				require.True(t, ok)

				assert.Contains(t, remoteLS, "TestEc2CopyFile")
			},
		})

	integration.ProgramTest(t, &test)
}

func TestLambdaInvoke(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "lambda-invoke"),
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
