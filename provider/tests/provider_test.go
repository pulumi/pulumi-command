package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blang/semver"
	"github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/property"

	command "github.com/pulumi/pulumi-command/provider/pkg/provider"
	"github.com/pulumi/pulumi-command/provider/pkg/provider/util/testutil"
	"github.com/pulumi/pulumi-command/provider/pkg/version"
)

// Property names and fixture values reused across the tests in this package.
const (
	createKey                 = "create"
	environmentKey            = "environment"
	stderrKey                 = "stderr"
	stdoutKey                 = "stdout"
	connectionKey             = "connection"
	hostKey                   = "host"
	addPreviousOutputInEnvKey = "addPreviousOutputInEnv"
	nameEnvVar                = "NAME"
	testResourceID            = "echo1234"
)

func provider(t *testing.T) integration.Server {
	v := semver.MustParse(version.Version)
	p, err := integration.NewServer(
		t.Context(),
		command.Name,
		v,
		integration.WithProvider(command.NewProvider()),
	)
	require.NoError(t, err)
	return p
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
	cmd := provider(t)
	urn := urn("local", "Command", "echo")
	property.New(property.Computed)
	computed := property.New(property.Computed)
	unknown := property.WithGoValue(computed, property.NewMap(nil))

	// Run a create against an in-memory provider, assert it succeeded, and return the
	// created property map.
	create := func(preview bool, env property.Value) property.Map {
		resp, err := cmd.Create(p.CreateRequest{
			Urn: urn,
			Properties: property.NewMap(map[string]property.Value{
				createKey:      property.New("echo hello, $NAME!"),
				environmentKey: env,
			}),
			DryRun: preview,
		})
		require.NoError(t, err)
		return resp.Properties
	}

	// The state that we expect a non-preview create to return.
	//
	// We use this as the final expect for create and the old state during update.
	createdState := property.NewMap(map[string]property.Value{
		createKey: property.New("echo hello, $NAME!"),
		stderrKey: property.New(""),
		stdoutKey: property.New("hello, world!"),
		environmentKey: property.New(property.NewMap(map[string]property.Value{
			nameEnvVar: property.New("world"),
		})),
	})

	// Run an update against an in-memory provider, assert it succeeded, and return
	// the new property map.
	update := func(preview bool, env property.Value) property.Map {
		resp, err := cmd.Update(p.UpdateRequest{
			ID:     testResourceID,
			Urn:    urn,
			DryRun: preview,
			State:  createdState,
			Inputs: property.NewMap(map[string]property.Value{
				createKey:      property.New("echo hello, $NAME!"),
				environmentKey: env,
			}),
		})
		require.NoError(t, err)
		return resp.Properties
	}

	t.Run("create-preview", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, property.NewMap(map[string]property.Value{
			createKey:      property.New("echo hello, $NAME!"),
			stderrKey:      computed,
			stdoutKey:      computed,
			environmentKey: unknown,
		}),
			create(true /* preview */, unknown))
	})

	t.Run("create-actual", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, createdState,
			create(false /* preview */, property.New(property.NewMap(map[string]property.Value{
				nameEnvVar: property.New("world"),
			}))))
	})

	t.Run("update-preview", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, property.NewMap(map[string]property.Value{
			createKey:      property.New("echo hello, $NAME!"),
			stderrKey:      computed,
			stdoutKey:      computed,
			environmentKey: unknown,
		}), update(true /* preview */, unknown))
	})
	t.Run("update-actual", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, property.NewMap(map[string]property.Value{
			createKey: property.New("echo hello, $NAME!"),
			environmentKey: property.New(property.NewMap(map[string]property.Value{
				nameEnvVar: property.New("Pulumi"),
			})),

			stderrKey: property.New(""),
			stdoutKey: property.New("hello, Pulumi!"),
		}),

			update(false /* preview */, property.New(property.NewMap(map[string]property.Value{
				nameEnvVar: property.New("Pulumi"),
			}))))
	})
}

func TestLocalCommandStdoutStderrFlag(t *testing.T) {
	cmd := provider(t)
	urn := urn("local", "Command", "echo")

	// Run a create against an in-memory provider, assert it succeeded, and return the
	// created property map.
	create := func() property.Map {
		resp, err := cmd.Create(p.CreateRequest{
			Urn: urn,
			Properties: property.NewMap(map[string]property.Value{
				createKey: property.New("echo std, $PULUMI_COMMAND_STDOUT"),
			}),
		})
		require.NoError(t, err)
		return resp.Properties
	}
	createdState := property.NewMap(map[string]property.Value{
		createKey: property.New("echo std, $PULUMI_COMMAND_STDOUT"),
		stderrKey: property.New(""),
		stdoutKey: property.New("std,"),
	})

	update := func(addPreviousOutputInEnv bool) property.Map {
		resp, err := cmd.Update(p.UpdateRequest{
			ID:    testResourceID,
			Urn:   urn,
			State: createdState,
			Inputs: property.NewMap(map[string]property.Value{
				createKey:                 property.New("echo std, $PULUMI_COMMAND_STDOUT"),
				addPreviousOutputInEnvKey: property.New(addPreviousOutputInEnv),
			}),
		})
		require.NoError(t, err)
		return resp.Properties
	}

	t.Run("create-actual", func(t *testing.T) {
		assert.Equal(t, createdState,
			create())
	})

	t.Run("update-actual-with-std", func(t *testing.T) {
		assert.Equal(t, property.NewMap(map[string]property.Value{
			createKey:                 property.New("echo std, $PULUMI_COMMAND_STDOUT"),
			stderrKey:                 property.New(""),
			stdoutKey:                 property.New("std, std,"),
			addPreviousOutputInEnvKey: property.New(true),
		}),

			update(true))
	})

	t.Run("update-actual-without-std", func(t *testing.T) {
		assert.Equal(t, property.NewMap(map[string]property.Value{
			createKey:                 property.New("echo std, $PULUMI_COMMAND_STDOUT"),
			stderrKey:                 property.New(""),
			stdoutKey:                 property.New("std,"),
			addPreviousOutputInEnvKey: property.New(false),
		}),

			update(false))
	})
}

func TestRemoteCommand(t *testing.T) {
	t.Parallel()

	t.Run("regress-256", func(t *testing.T) {
		resp, err := provider(t).Create(p.CreateRequest{
			Urn:    urn("remote", "Command", "check"),
			DryRun: true,
			Properties: property.NewMap(map[string]property.Value{
				createKey: property.New("<create command>"),
				connectionKey: property.New(property.NewMap(map[string]property.Value{
					hostKey: property.New("<host port>"),
				})).WithSecret(true),
			}),
		})
		require.NoError(t, err)

		for _, key := range []string{stdoutKey, stderrKey} {
			p := resp.Properties.Get(key)
			assert.True(t, p.HasComputed())
			assert.False(t, p.HasSecrets())
			// assert.False(t, p.IsSecret() || (p.IsOutput() && p.OutputValue().Secret))
		}
	})
}

func TestRemoteCommandStdoutStderrFlag(t *testing.T) {
	// Start a local SSH server that writes the PULUMI_COMMAND_STDOUT environment variable
	// on the format "PULUMI_COMMAND_STDOUT=<value>" to the client using stdout.
	const (
		createCommand = "arbitrary create command"
	)

	sshServer := testutil.NewTestSSHServer(t, func(session ssh.Session) {
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

	cmd := provider(t)
	urn := urn("remote", "Command", "dial")

	// Run a create against an in-memory provider, assert it succeeded, and return the created property map.
	connection := property.New(property.NewMap(map[string]property.Value{
		hostKey:          property.New(sshServer.Host),
		"port":           property.New(float64(sshServer.Port)),
		"user":           property.New("arbitrary-user"), // unused but prevents nil panic
		"perDialTimeout": property.New(1.0),
	}))

	// unused but prevents nil panic

	// The state that we expect a non-preview create to return.
	//
	// We use this as the final expect for create and the old state during update.
	initialState := property.NewMap(map[string]property.Value{
		connectionKey:             connection,
		createKey:                 property.New(createCommand),
		stderrKey:                 property.New(""),
		stdoutKey:                 property.New("Response{}"),
		addPreviousOutputInEnvKey: property.New(true),
	})

	t.Run(createKey, func(t *testing.T) {
		createResponse, err := cmd.Create(p.CreateRequest{
			Urn: urn,
			Properties: property.NewMap(map[string]property.Value{
				connectionKey:             connection,
				createKey:                 property.New(createCommand),
				addPreviousOutputInEnvKey: property.New(true),
			}),
		})
		require.NoError(t, err)
		require.Equal(t, initialState, createResponse.Properties)
	})

	// Run an update against an in-memory provider, assert it succeeded, and return
	// the new property map.
	update := func(addPreviousOutputInEnv bool) property.Map {
		resp, err := cmd.Update(p.UpdateRequest{
			ID:    testResourceID,
			Urn:   urn,
			State: initialState,
			Inputs: property.NewMap(map[string]property.Value{
				connectionKey:             connection,
				createKey:                 property.New(createCommand),
				addPreviousOutputInEnvKey: property.New(addPreviousOutputInEnv),
			}),
		})
		require.NoError(t, err)
		return resp.Properties
	}

	t.Run("update-actual-with-std", func(t *testing.T) {
		assert.Equal(t, property.NewMap(map[string]property.Value{
			connectionKey: connection,
			createKey:     property.New(createCommand),
			stderrKey:     property.New(""),
			// Running with addPreviousOutputInEnv=true sets the environment variable:
			stdoutKey:                 property.New("Response{PULUMI_COMMAND_STDOUT=Response{}}"),
			addPreviousOutputInEnvKey: property.New(true),
		}),

			update(true))
	})

	t.Run("update-actual-without-std", func(t *testing.T) {
		assert.Equal(t, property.NewMap(map[string]property.Value{
			connectionKey: connection,
			createKey:     property.New(createCommand),
			stderrKey:     property.New(""),
			// Running without addPreviousOutputInEnv does not set the environment variable:
			stdoutKey:                 property.New("Response{}"),
			addPreviousOutputInEnvKey: property.New(false),
		}),

			update(false))
	})
}

// Ensure that we correctly apply defaults to `connection.port`.
//
// User issue is https://github.com/pulumi/pulumi-command/issues/248.
func TestRegress248(t *testing.T) {
	t.Parallel()
	resp, err := provider(t).Check(p.CheckRequest{
		Urn: urn("remote", "Command", "check"),
		Inputs: property.NewMap(map[string]property.Value{
			createKey: property.New("<create command>"),
			connectionKey: property.New(property.NewMap(map[string]property.Value{
				hostKey: property.New("<required value>"),
			})),
		}),
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Failures)
	assert.Equal(t, property.NewMap(map[string]property.Value{
		createKey: property.New("<create command>"),
		connectionKey: property.New(property.NewMap(map[string]property.Value{
			hostKey:          property.New("<required value>"),
			"port":           property.New(22.0),
			"user":           property.New("root"),
			"dialErrorLimit": property.New(10.0),
			"perDialTimeout": property.New(15.0),
		})).WithSecret(true),
	}),

		resp.Inputs)
}

func TestLocalRun(t *testing.T) {
	t.Parallel()

	resp, err := provider(t).Invoke(p.InvokeRequest{
		Token: "command:local:run",
		Args: property.NewMap(map[string]property.Value{
			"command": property.New(`echo "Hello, World!"`),
		}),
	})
	require.NoError(t, err)
	assert.Equal(t, property.NewMap(map[string]property.Value{
		"command": property.New(`echo "Hello, World!"`),
		stderrKey: property.New(""),
		stdoutKey: property.New("Hello, World!"),
	}), resp.Return)
}
