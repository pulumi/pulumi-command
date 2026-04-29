// Copyright 2016-2021, Pulumi Corporation.  All rights reserved.
//go:build nodejs || all
// +build nodejs all

package examples

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/pulumi/providertest"
	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/assertpreview"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/pulumi-command/examples/sshfixture"
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

// remoteBaseConfig returns the Config and Secrets blocks describing a single
// SSH server fixture.
func remoteBaseConfig(s sshfixture.Server) (cfg, secrets map[string]string) {
	cfg = map[string]string{
		"host": s.Host,
		"port": strconv.Itoa(s.Port),
		"user": sshfixture.User,
	}
	secrets = map[string]string{
		"privateKeyBase64": base64.StdEncoding.EncodeToString([]byte(s.PrivateKeyPEM)),
	}
	return cfg, secrets
}

func TestRemoteNodejs(t *testing.T) {
	server := sshfixture.New(t)
	cfg, secrets := remoteBaseConfig(server)

	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:     filepath.Join(getCwd(t), "remote-nodejs"),
			Config:  cfg,
			Secrets: secrets,
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
					sig, sigOk := sigKey.(string)
					if !sigOk || sig != resource.SecretSig {
						return false
					}
					_, cOk := m["ciphertext"].(string)
					return cOk
				}
				assert.Truef(t, isEncrypted(stack.Outputs["connectionSecret"]),
					"connectionSecret value should be encrypted")
				assert.Equal(t, "micro", fmt.Sprintf("%v", stack.Outputs["confirmSize"]))
			},
		})

	integration.ProgramTest(t, &test)
}

func TestRemoteProxyNodejs(t *testing.T) {
	proxy := sshfixture.NewProxy(t)

	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "remote-proxy-nodejs"),
			Config: map[string]string{
				"proxyHost":  proxy.Proxy.Host,
				"proxyPort":  strconv.Itoa(proxy.Proxy.Port),
				"targetHost": proxy.Target.Host,
				"targetPort": strconv.Itoa(proxy.Target.Port),
				"user":       sshfixture.User,
			},
			Secrets: map[string]string{
				"privateKeyBase64": base64.StdEncoding.EncodeToString([]byte(proxy.PrivateKeyPEM)),
			},
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				assert.Equal(t, "micro", fmt.Sprintf("%v", stack.Outputs["confirmSize"]))
			},
		})

	integration.ProgramTest(t, &test)
}

func TestDirCopyNodejs(t *testing.T) {
	server := sshfixture.New(t)
	cfg, secrets := remoteBaseConfig(server)

	const dest = "/tmp/dir-copy-nodejs"
	const extrasDest = dest + "-extras"
	cfg["destDir"] = dest
	basePath := filepath.Join(getCwd(t), "dir-copy-nodejs")

	expectedExtras := extrasDest + "\n" +
		extrasDest + "/asset-archive\n" +
		extrasDest + "/asset-archive/greeting.txt\n" +
		extrasDest + "/asset-archive/nested\n" +
		extrasDest + "/asset-archive/nested/answer.txt\n" +
		extrasDest + "/remote-archive.tar.gz\n" +
		extrasDest + "/remote-asset.txt\n" +
		extrasDest + "/string-asset.txt"

	assertExtras := func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
		assert.Equal(t, expectedExtras, stringOutput(t, stack, "lsExtras"))
	}

	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:     basePath,
			Config:  cfg,
			Secrets: secrets,
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				remoteLS := stringOutput(t, stack, "lsRemote")
				assert.Equal(t,
					dest+"\n"+
						dest+"/file1\n"+
						dest+"/one\n"+
						dest+"/one/file2\n"+
						dest+"/one/two\n"+
						dest+"/one/two/file3",
					remoteLS)

				assertExtras(t, stack)
			},
			EditDirs: []integration.EditDir{
				{
					Dir:             filepath.Join(basePath, "step2"),
					Additive:        true,
					ExpectNoChanges: false,
					ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
						remoteLS := stringOutput(t, stack, "lsRemote")
						assert.Equal(t,
							dest+"\n"+
								dest+"/file1\n"+
								dest+"/newfile\n"+
								dest+"/one\n"+
								dest+"/one/file2\n"+
								dest+"/one/two\n"+
								dest+"/one/two/file3",
							remoteLS)

						assertExtras(t, stack)
					},
				},
			},
		})

	integration.ProgramTest(t, &test)
}

func TestCopyFileNodejs(t *testing.T) {
	server := sshfixture.New(t)
	cfg, secrets := remoteBaseConfig(server)

	const dest = "/tmp/TestCopyFileNodejs"
	cfg["destDir"] = dest

	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:     filepath.Join(getCwd(t), "copyfile-nodejs"),
			Config:  cfg,
			Secrets: secrets,
			ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
				remoteLS := stringOutput(t, stack, "lsRemote")
				assert.Equal(t, dest, remoteLS)
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

func TestUpgradeLocalCommand(t *testing.T) {
	t.Parallel()

	dir := fmt.Sprintf("./stdin")

	test := pulumitest.NewPulumiTest(t, dir,
		opttest.YarnLink("@pulumi/command"),
		opttest.LocalProviderPath("command", filepath.Join(dir, "../../", "bin")),
	)
	result := providertest.PreviewProviderUpgrade(t, test, "command", "0.11.1")
	assertpreview.HasNoChanges(t, result)
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

func stringOutput(t *testing.T, stack integration.RuntimeValidationStackInfo, key string) string {
	t.Helper()
	out, ok := stack.Outputs[key]
	require.Truef(t, ok, "missing output %q", key)
	s, ok := out.(string)
	require.Truef(t, ok, "output %q is not a string: %T", key, out)
	return s
}
