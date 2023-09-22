package tests

import (
	"fmt"
	"testing"

	"github.com/blang/semver"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	command "github.com/pulumi/pulumi-command/provider/pkg/provider"
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
	n := tokens.QName(name)
	return resource.NewURN("test", "command", "",
		tokens.NewTypeToken(
			tokens.NewModuleToken(command.Name, m),
			r),
		n)
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

// Ensure that we correctly apply apply defaults to `connection.port`.
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
		"connection": resource.NewObjectProperty(resource.PropertyMap{
			"host":           pString("<required value>"),
			"port":           pNumber(22),
			"user":           pString("root"),
			"dialErrorLimit": pNumber(10),
			"perDialTimeout": pNumber(15),
		}),
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
		"command": pString(`echo "Hello, World!"`),
		"stderr":  pString(""),
		"stdout":  pString("Hello, World!"),
	}, resp.Return)
}
