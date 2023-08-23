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
		Element: resource.NewStringProperty(""),
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
			"create": resource.PropertyValue{V: "echo hello, $NAME!"},
			"stderr": c(resource.PropertyValue{V: ""}),
			"stdout": c(resource.PropertyValue{V: ""}),
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
			"create": resource.PropertyValue{V: "echo hello, $NAME!"},
			"stderr": c(resource.PropertyValue{V: ""}),
			"stdout": c(resource.PropertyValue{V: "hello, world!"}),
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
