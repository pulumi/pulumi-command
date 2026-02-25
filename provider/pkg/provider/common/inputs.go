package common //nolint:revive

import "github.com/pulumi/pulumi-go-provider/infer"

// ResourceInputs are inputs common to resource CRUD operations.
type ResourceInputs struct {
	// The field tags are used to provide metadata on the schema representation.
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	Triggers *[]any `pulumi:"triggers,optional" provider:"replaceOnChanges"`

	Create *string `pulumi:"create,optional"`
	Update *string `pulumi:"update,optional"`
	Delete *string `pulumi:"delete,optional"`
}

// Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *ResourceInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Triggers, "The resource will be updated (or replaced) if any of these values change.\n\n"+
		"The trigger values can be of any type.\n\n"+
		"If the `update` command was provided the resource will be updated, otherwise it will be replaced "+
		"using the `create` command.\n\n"+
		"Please see the resource documentation for examples.",
	)
	a.Describe(&c.Create, "The command to run once on resource creation.\n\n"+
		"If an `update` command isn't provided, then `create` will also be run when the resource's "+
		"inputs are modified.\n\n"+
		"Note that this command will not be executed if the resource has already been created and "+
		"its inputs are unchanged.\n\n"+
		"Use `local.runOutput` if you need to run a command on every execution of your program.",
	)
	a.Describe(&c.Delete, "The command to run on resource delettion.\n\n"+
		"The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to "+
		"the stdout and stderr properties of the Command resource from previous create or update steps.",
	)
	a.Describe(&c.Update, "The command to run when the resource is updated.\n\n"+
		"If empty, the create command will be executed instead.\n\n"+
		"Note that this command will not run if the resource's inputs are unchanged.\n\n"+
		"Use `local.runOutput` if you need to run a command on every execution of your program.\n\n"+
		"The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to "+
		"the `stdout` and `stderr` properties of the Command resource from previous create or update steps.",
	)
}
