package common

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
	a.Describe(&c.Triggers, `Trigger a resource replacement on changes to any of these values. The
trigger values can be of any type. If a value is different in the current update compared to the
previous update, the resource will be replaced, i.e., the "create" command will be re-run.

Here are some examples of triggers in TypeScript. Changes to any of the four will cause `+"`cmd`"+
		` to be re-run. However, note that for `+"`fileAsset`"+` it's the variable itself that is
the trigger, not the contents of index.ts, since triggers are simply opaque values.

`+"```typescript"+
		`import * as local from "@pulumi/command/local";
import * as random from "@pulumi/random";
import { asset } from "@pulumi/pulumi";
import * as path from "path";

const str = "foo";

const rand = new random.RandomString("rand", { length: 5 });

const fileAsset = new asset.FileAsset(path.join(__dirname, "index.ts"));

const localFile = new local.Command("localFile", {
    create: `+"`touch foo.txt`"+`,
    archivePaths: ["*.txt"]
});

const cmd = new local.Command("pwd", {
    create: `+"`echo create > op.txt`"+`,
    delete: `+"`echo delete >> op.txt`"+`,
    triggers: [str, rand.result, fileAsset, localFile.archive],
});`+
		"```")

	a.Describe(&c.Create, "The command to run on create.")
	a.Describe(&c.Delete, `The command to run on delete. The environment variables PULUMI_COMMAND_STDOUT
and PULUMI_COMMAND_STDERR are set to the stdout and stderr properties of the
Command resource from previous create or update steps.`)
	a.Describe(&c.Update, `The command to run on update, if empty, create will 
run again. The environment variables PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR 
are set to the stdout and stderr properties of the Command resource from previous 
create or update steps.`)
}
