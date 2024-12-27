package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blang/semver"
	"github.com/gliderlabs/ssh"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	command "github.com/pulumi/pulumi-command/provider/pkg/provider"
	"github.com/pulumi/pulumi-command/provider/pkg/provider/util/testutil"
	"github.com/pulumi/pulumi-command/provider/pkg/version"
)

func provider() integration.Server {
	v := semver.MustParse(version.Version)
	return integration.NewServer(command.Name, v, command.NewProvider())
}

func urn(mod, res, name string) resource.URN {
	m := tokens.ModuleName(mod)
	r := tokens.TypeName(res)
	if !tokens.IsQName(name) {
		panic(fmt.Sprintf("invalid resource name: %q", name))
	}
	return resource.NewURN("test", "command", "",
		tokens.NewTypeToken(
			tokens.NewModuleToken(command.Name, m),
			r),
		name)
}

func TestLocalCommand(t *testing.T) {
	t.Parallel()
	cmd := provider()
	urn := urn("local", "Command", "echo")
	unknown := resource.NewOutputProperty(resource.Output{
		Element: resource.NewObjectProperty(resource.PropertyMap{}),
		Known:   false,
	})
	c := resource.MakeComputed

	// Run a create against an in-memory provider, assert it succeeded, and return the
	// created property map.
	create := func(preview bool, env resource.PropertyValue) resource.PropertyMap {
		resp, err := cmd.Create(p.CreateRequest{
			Urn: urn,
			Properties: resource.PropertyMap{
				"create":      resource.NewStringProperty("echo hello, $NAME!"),
				"environment": env,
			},
			Preview: preview,
		})
		require.NoError(t, err)
		return resp.Properties
	}

	// The state that we expect a non-preview create to return.
	//
	// We use this as the final expect for create and the old state during update.
	createdState := resource.PropertyMap{
		"create": resource.PropertyValue{V: "echo hello, $NAME!"},
		"stderr": resource.PropertyValue{V: ""},
		"stdout": resource.PropertyValue{V: "hello, world!"},
		"environment": resource.NewObjectProperty(resource.PropertyMap{
			"NAME": resource.NewStringProperty("world"),
		}),
	}

	// Run an update against an in-memory provider, assert it succeeded, and return
	// the new property map.
	update := func(preview bool, env resource.PropertyValue) resource.PropertyMap {
		resp, err := cmd.Update(p.UpdateRequest{
			ID:      "echo1234",
			Urn:     urn,
			Preview: preview,
			Olds:    createdState.Copy(),
			News: resource.PropertyMap{
				"create":      resource.NewStringProperty("echo hello, $NAME!"),
				"environment": env,
			},
		})
		require.NoError(t, err)
		return resp.Properties
	}

	t.Run("create-preview", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, resource.PropertyMap{
			"create":      resource.PropertyValue{V: "echo hello, $NAME!"},
			"stderr":      c(resource.PropertyValue{V: ""}),
			"stdout":      c(resource.PropertyValue{V: ""}),
			"environment": unknown,
		},
			create(true /* preview */, unknown))
	})

	t.Run("create-actual", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, createdState,
			create(false /* preview */, resource.NewObjectProperty(resource.PropertyMap{
				"NAME": resource.NewStringProperty("world"),
			})))
	})

	t.Run("update-preview", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, resource.PropertyMap{
			"create":      resource.PropertyValue{V: "echo hello, $NAME!"},
			"stderr":      c(resource.PropertyValue{V: ""}),
			"stdout":      c(resource.PropertyValue{V: "hello, world!"}),
			"environment": unknown,
		}, update(true /* preview */, unknown))
	})
	t.Run("update-actual", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, resource.PropertyMap{
			"create": resource.PropertyValue{V: "echo hello, $NAME!"},
			"environment": resource.NewObjectProperty(resource.PropertyMap{
				"NAME": resource.NewStringProperty("Pulumi"),
			}),
			"stderr": resource.PropertyValue{V: ""},
			"stdout": resource.PropertyValue{V: "hello, Pulumi!"},
		}, update(false /* preview */, resource.NewObjectProperty(resource.PropertyMap{
			"NAME": resource.NewStringProperty("Pulumi"),
		})))
	})
}

func TestLocalCommandStdoutStderrFlag(t *testing.T) {
	cmd := provider()
	urn := urn("local", "Command", "echo")

	// Run a create against an in-memory provider, assert it succeeded, and return the
	// created property map.
	create := func() resource.PropertyMap {
		resp, err := cmd.Create(p.CreateRequest{
			Urn: urn,
			Properties: resource.PropertyMap{
				"create": resource.NewStringProperty("echo std, $PULUMI_COMMAND_STDOUT"),
			},
		})
		require.NoError(t, err)
		return resp.Properties
	}

	// The state that we expect a non-preview create to return.
	//
	// We use this as the final expect for create and the old state during update.
	createdState := resource.PropertyMap{
		"create": resource.PropertyValue{V: "echo std, $PULUMI_COMMAND_STDOUT"},
		"stderr": resource.PropertyValue{V: ""},
		"stdout": resource.PropertyValue{V: "std,"},
	}

	// Run an update against an in-memory provider, assert it succeeded, and return
	// the new property map.
	update := func(addPreviousOutputInEnv bool) resource.PropertyMap {
		resp, err := cmd.Update(p.UpdateRequest{
			ID:   "echo1234",
			Urn:  urn,
			Olds: createdState.Copy(),
			News: resource.PropertyMap{
				"create":                 resource.NewStringProperty("echo std, $PULUMI_COMMAND_STDOUT"),
				"addPreviousOutputInEnv": resource.NewBoolProperty(addPreviousOutputInEnv),
			},
		})
		require.NoError(t, err)
		return resp.Properties
	}

	t.Run("create-actual", func(t *testing.T) {
		assert.Equal(t, createdState,
			create())
	})

	t.Run("update-actual-with-std", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"create":                 resource.PropertyValue{V: "echo std, $PULUMI_COMMAND_STDOUT"},
			"stderr":                 resource.PropertyValue{V: ""},
			"stdout":                 resource.PropertyValue{V: "std, std,"},
			"addPreviousOutputInEnv": resource.PropertyValue{V: true},
		}, update(true))
	})

	t.Run("update-actual-without-std", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"create":                 resource.PropertyValue{V: "echo std, $PULUMI_COMMAND_STDOUT"},
			"stderr":                 resource.PropertyValue{V: ""},
			"stdout":                 resource.PropertyValue{V: "std,"},
			"addPreviousOutputInEnv": resource.PropertyValue{V: false},
		}, update(false))
	})

}

func TestRemoteCommand(t *testing.T) {
	t.Parallel()
	pString := resource.NewStringProperty
	sec := resource.MakeSecret

	t.Run("regress-256", func(t *testing.T) {
		resp, err := provider().Create(p.CreateRequest{
			Urn:     urn("remote", "Command", "check"),
			Preview: true,
			Properties: resource.PropertyMap{
				"create": pString("<create command>"),
				"connection": sec(resource.NewObjectProperty(resource.PropertyMap{
					"host": pString("<host port>"),
				})),
			}})
		require.NoError(t, err)

		for _, v := range []resource.PropertyKey{"stdout", "stderr"} {
			p := resp.Properties[v]
			assert.True(t, p.ContainsUnknowns())
			assert.False(t, p.IsSecret() || (p.IsOutput() && p.OutputValue().Secret))
		}
	})
}

func TestRemoteCommandStdoutStderrFlag(t *testing.T) {
	// Start a local SSH server that writes the PULUMI_COMMAND_STDOUT environment variable
	// on the format "PULUMI_COMMAND_STDOUT=<value>" to the client using stdout.
	const (
		createCommand = "arbitrary create command"
	)

	sshServer := testutil.NewTestSshServer(t, func(session ssh.Session) {
		// Find the PULUMI_COMMAND_STDOUT environment variable
		var envVar string
		for _, v := range session.Environ() {
			if strings.HasPrefix(v, "PULUMI_COMMAND_STDOUT=") {
				envVar = v
				break
			}
		}

		response := fmt.Sprintf("Response{%s}", envVar)
		_, err := session.Write([]byte(response))
		require.NoErrorf(t, err, "session.Write(%s)", response)
	})

	cmd := provider()
	urn := urn("remote", "Command", "dial")

	// Run a create against an in-memory provider, assert it succeeded, and return the created property map.
	connection := resource.NewObjectProperty(resource.PropertyMap{
		"host":           resource.NewStringProperty(sshServer.Host),
		"port":           resource.NewNumberProperty(float64(sshServer.Port)),
		"user":           resource.NewStringProperty("arbitrary-user"), // unused but prevents nil panic
		"perDialTimeout": resource.NewNumberProperty(1),                // unused but prevents nil panic
	})

	// The state that we expect a non-preview create to return.
	//
	// We use this as the final expect for create and the old state during update.
	initialState := resource.PropertyMap{
		"connection":             connection,
		"create":                 resource.PropertyValue{V: createCommand},
		"stderr":                 resource.PropertyValue{V: ""},
		"stdout":                 resource.PropertyValue{V: "Response{}"},
		"addPreviousOutputInEnv": resource.NewBoolProperty(true),
	}

	t.Run("create", func(t *testing.T) {
		createResponse, err := cmd.Create(p.CreateRequest{
			Urn: urn,
			Properties: resource.PropertyMap{
				"connection":             connection,
				"create":                 resource.NewStringProperty(createCommand),
				"addPreviousOutputInEnv": resource.NewBoolProperty(true),
			},
		})
		require.NoError(t, err)
		require.Equal(t, initialState, createResponse.Properties)
	})

	// Run an update against an in-memory provider, assert it succeeded, and return
	// the new property map.
	update := func(addPreviousOutputInEnv bool) resource.PropertyMap {
		resp, err := cmd.Update(p.UpdateRequest{
			ID:   "echo1234",
			Urn:  urn,
			Olds: initialState.Copy(),
			News: resource.PropertyMap{
				"connection":             connection,
				"create":                 resource.NewStringProperty(createCommand),
				"addPreviousOutputInEnv": resource.NewBoolProperty(addPreviousOutputInEnv),
			},
		})
		require.NoError(t, err)
		return resp.Properties
	}

	t.Run("update-actual-with-std", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"connection": connection,
			"create":     resource.PropertyValue{V: createCommand},
			"stderr":     resource.PropertyValue{V: ""},
			// Running with addPreviousOutputInEnv=true sets the environment variable:
			"stdout":                 resource.PropertyValue{V: "Response{PULUMI_COMMAND_STDOUT=Response{}}"},
			"addPreviousOutputInEnv": resource.PropertyValue{V: true},
		}, update(true))
	})

	t.Run("update-actual-without-std", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"connection": connection,
			"create":     resource.PropertyValue{V: createCommand},
			"stderr":     resource.PropertyValue{V: ""},
			// Running without addPreviousOutputInEnv does not set the environment variable:
			"stdout":                 resource.PropertyValue{V: "Response{}"},
			"addPreviousOutputInEnv": resource.PropertyValue{V: false},
		}, update(false))
	})

}

// Ensure that we correctly apply defaults to `connection.port`.
//
// User issue is https://github.com/pulumi/pulumi-command/issues/248.
func TestRegress248(t *testing.T) {
	t.Parallel()
	type pMap = resource.PropertyMap
	pString := resource.NewStringProperty
	pNumber := resource.NewNumberProperty
	resp, err := provider().Check(p.CheckRequest{
		Urn: urn("remote", "Command", "check"),
		News: resource.PropertyMap{
			"create": pString("<create command>"),
			"connection": resource.NewObjectProperty(pMap{
				"host": pString("<required value>"),
			}),
		},
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Failures)
	assert.Equal(t, resource.PropertyMap{
		"create": pString("<create command>"),
		"connection": resource.MakeSecret(resource.NewObjectProperty(resource.PropertyMap{
			"host":           pString("<required value>"),
			"port":           pNumber(22),
			"user":           pString("root"),
			"dialErrorLimit": pNumber(10),
			"perDialTimeout": pNumber(15),
		})),
		"addPreviousOutputInEnv": resource.NewBoolProperty(true),
	}, resp.Inputs)
}

func TestLocalRun(t *testing.T) {
	t.Parallel()

	type pMap = resource.PropertyMap
	pString := resource.NewStringProperty

	resp, err := provider().Invoke(p.InvokeRequest{
		Token: "command:local:run",
		Args: pMap{
			"command": pString(`echo "Hello, World!"`),
		},
	})
	require.NoError(t, err)
	assert.Equal(t, pMap{
		"command":                pString(`echo "Hello, World!"`),
		"stderr":                 pString(""),
		"stdout":                 pString("Hello, World!"),
		"addPreviousOutputInEnv": resource.NewProperty(true),
	}, resp.Return)
}
