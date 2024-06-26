// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package remote

import (
	"context"
	"reflect"

	"errors"
	"github.com/pulumi/pulumi-command/sdk/go/command/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Copy an Asset or Archive to a remote host.
//
// ## Example usage
//
// This example copies a local directory to a remote host via SSH. For brevity, the remote server is assumed to exist, but it could also be provisioned in the same Pulumi program.
//
// ```go
// package main
//
// import (
//
//	"fmt"
//
//	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			cfg := config.New(ctx, "")
//			serverPublicIp := cfg.Require("serverPublicIp")
//			userName := cfg.Require("userName")
//			privateKey := cfg.Require("privateKey")
//			payload := cfg.Require("payload")
//			destDir := cfg.Require("destDir")
//
//			archive := pulumi.NewFileArchive(payload)
//
//			conn := remote.ConnectionArgs{
//				Host:       pulumi.String(serverPublicIp),
//				User:       pulumi.String(userName),
//				PrivateKey: pulumi.String(privateKey),
//			}
//
//			copy, err := remote.NewCopyToRemote(ctx, "copy", &remote.CopyToRemoteArgs{
//				Connection: conn,
//				Source:     archive,
//			})
//			if err != nil {
//				return err
//			}
//
//			find, err := remote.NewCommand(ctx, "find", &remote.CommandArgs{
//				Connection: conn,
//				Create:     pulumi.String(fmt.Sprintf("find %v/%v | sort", destDir, payload)),
//				Triggers: pulumi.Array{
//					archive,
//				},
//			}, pulumi.DependsOn([]pulumi.Resource{
//				copy,
//			}))
//			if err != nil {
//				return err
//			}
//
//			ctx.Export("remoteContents", find.Stdout)
//			return nil
//		})
//	}
//
// ```
type CopyToRemote struct {
	pulumi.CustomResourceState

	// The parameters with which to connect to the remote host.
	Connection ConnectionOutput `pulumi:"connection"`
	// The destination path on the remote host. The last element of the path will be created if it doesn't exist but it's an error when additional elements don't exist. When the remote path is an existing directory, the source file or directory will be copied into that directory. When the source is a file and the remote path is an existing file, that file will be overwritten. When the source is a directory and the remote path an existing file, the copy will fail.
	RemotePath pulumi.StringOutput `pulumi:"remotePath"`
	// An [asset or an archive](https://www.pulumi.com/docs/concepts/assets-archives/) to upload as the source of the copy. It must be path-based, i.e., be a `FileAsset` or a `FileArchive`. The item will be copied as-is; archives like .tgz will not be unpacked. Directories are copied recursively, overwriting existing files.
	Source pulumi.AssetOrArchiveOutput `pulumi:"source"`
	// Trigger replacements on changes to this input.
	Triggers pulumi.ArrayOutput `pulumi:"triggers"`
}

// NewCopyToRemote registers a new resource with the given unique name, arguments, and options.
func NewCopyToRemote(ctx *pulumi.Context,
	name string, args *CopyToRemoteArgs, opts ...pulumi.ResourceOption) (*CopyToRemote, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Connection == nil {
		return nil, errors.New("invalid value for required argument 'Connection'")
	}
	if args.RemotePath == nil {
		return nil, errors.New("invalid value for required argument 'RemotePath'")
	}
	if args.Source == nil {
		return nil, errors.New("invalid value for required argument 'Source'")
	}
	args.Connection = args.Connection.ToConnectionOutput().ApplyT(func(v Connection) Connection { return *v.Defaults() }).(ConnectionOutput)
	if args.Connection != nil {
		args.Connection = pulumi.ToSecret(args.Connection).(ConnectionInput)
	}
	secrets := pulumi.AdditionalSecretOutputs([]string{
		"connection",
	})
	opts = append(opts, secrets)
	replaceOnChanges := pulumi.ReplaceOnChanges([]string{
		"triggers[*]",
	})
	opts = append(opts, replaceOnChanges)
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource CopyToRemote
	err := ctx.RegisterResource("command:remote:CopyToRemote", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetCopyToRemote gets an existing CopyToRemote resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetCopyToRemote(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *CopyToRemoteState, opts ...pulumi.ResourceOption) (*CopyToRemote, error) {
	var resource CopyToRemote
	err := ctx.ReadResource("command:remote:CopyToRemote", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering CopyToRemote resources.
type copyToRemoteState struct {
}

type CopyToRemoteState struct {
}

func (CopyToRemoteState) ElementType() reflect.Type {
	return reflect.TypeOf((*copyToRemoteState)(nil)).Elem()
}

type copyToRemoteArgs struct {
	// The parameters with which to connect to the remote host.
	Connection Connection `pulumi:"connection"`
	// The destination path on the remote host. The last element of the path will be created if it doesn't exist but it's an error when additional elements don't exist. When the remote path is an existing directory, the source file or directory will be copied into that directory. When the source is a file and the remote path is an existing file, that file will be overwritten. When the source is a directory and the remote path an existing file, the copy will fail.
	RemotePath string `pulumi:"remotePath"`
	// An [asset or an archive](https://www.pulumi.com/docs/concepts/assets-archives/) to upload as the source of the copy. It must be path-based, i.e., be a `FileAsset` or a `FileArchive`. The item will be copied as-is; archives like .tgz will not be unpacked. Directories are copied recursively, overwriting existing files.
	Source pulumi.AssetOrArchive `pulumi:"source"`
	// Trigger replacements on changes to this input.
	Triggers []interface{} `pulumi:"triggers"`
}

// The set of arguments for constructing a CopyToRemote resource.
type CopyToRemoteArgs struct {
	// The parameters with which to connect to the remote host.
	Connection ConnectionInput
	// The destination path on the remote host. The last element of the path will be created if it doesn't exist but it's an error when additional elements don't exist. When the remote path is an existing directory, the source file or directory will be copied into that directory. When the source is a file and the remote path is an existing file, that file will be overwritten. When the source is a directory and the remote path an existing file, the copy will fail.
	RemotePath pulumi.StringInput
	// An [asset or an archive](https://www.pulumi.com/docs/concepts/assets-archives/) to upload as the source of the copy. It must be path-based, i.e., be a `FileAsset` or a `FileArchive`. The item will be copied as-is; archives like .tgz will not be unpacked. Directories are copied recursively, overwriting existing files.
	Source pulumi.AssetOrArchiveInput
	// Trigger replacements on changes to this input.
	Triggers pulumi.ArrayInput
}

func (CopyToRemoteArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*copyToRemoteArgs)(nil)).Elem()
}

type CopyToRemoteInput interface {
	pulumi.Input

	ToCopyToRemoteOutput() CopyToRemoteOutput
	ToCopyToRemoteOutputWithContext(ctx context.Context) CopyToRemoteOutput
}

func (*CopyToRemote) ElementType() reflect.Type {
	return reflect.TypeOf((**CopyToRemote)(nil)).Elem()
}

func (i *CopyToRemote) ToCopyToRemoteOutput() CopyToRemoteOutput {
	return i.ToCopyToRemoteOutputWithContext(context.Background())
}

func (i *CopyToRemote) ToCopyToRemoteOutputWithContext(ctx context.Context) CopyToRemoteOutput {
	return pulumi.ToOutputWithContext(ctx, i).(CopyToRemoteOutput)
}

// CopyToRemoteArrayInput is an input type that accepts CopyToRemoteArray and CopyToRemoteArrayOutput values.
// You can construct a concrete instance of `CopyToRemoteArrayInput` via:
//
//	CopyToRemoteArray{ CopyToRemoteArgs{...} }
type CopyToRemoteArrayInput interface {
	pulumi.Input

	ToCopyToRemoteArrayOutput() CopyToRemoteArrayOutput
	ToCopyToRemoteArrayOutputWithContext(context.Context) CopyToRemoteArrayOutput
}

type CopyToRemoteArray []CopyToRemoteInput

func (CopyToRemoteArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*CopyToRemote)(nil)).Elem()
}

func (i CopyToRemoteArray) ToCopyToRemoteArrayOutput() CopyToRemoteArrayOutput {
	return i.ToCopyToRemoteArrayOutputWithContext(context.Background())
}

func (i CopyToRemoteArray) ToCopyToRemoteArrayOutputWithContext(ctx context.Context) CopyToRemoteArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(CopyToRemoteArrayOutput)
}

// CopyToRemoteMapInput is an input type that accepts CopyToRemoteMap and CopyToRemoteMapOutput values.
// You can construct a concrete instance of `CopyToRemoteMapInput` via:
//
//	CopyToRemoteMap{ "key": CopyToRemoteArgs{...} }
type CopyToRemoteMapInput interface {
	pulumi.Input

	ToCopyToRemoteMapOutput() CopyToRemoteMapOutput
	ToCopyToRemoteMapOutputWithContext(context.Context) CopyToRemoteMapOutput
}

type CopyToRemoteMap map[string]CopyToRemoteInput

func (CopyToRemoteMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*CopyToRemote)(nil)).Elem()
}

func (i CopyToRemoteMap) ToCopyToRemoteMapOutput() CopyToRemoteMapOutput {
	return i.ToCopyToRemoteMapOutputWithContext(context.Background())
}

func (i CopyToRemoteMap) ToCopyToRemoteMapOutputWithContext(ctx context.Context) CopyToRemoteMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(CopyToRemoteMapOutput)
}

type CopyToRemoteOutput struct{ *pulumi.OutputState }

func (CopyToRemoteOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**CopyToRemote)(nil)).Elem()
}

func (o CopyToRemoteOutput) ToCopyToRemoteOutput() CopyToRemoteOutput {
	return o
}

func (o CopyToRemoteOutput) ToCopyToRemoteOutputWithContext(ctx context.Context) CopyToRemoteOutput {
	return o
}

// The parameters with which to connect to the remote host.
func (o CopyToRemoteOutput) Connection() ConnectionOutput {
	return o.ApplyT(func(v *CopyToRemote) ConnectionOutput { return v.Connection }).(ConnectionOutput)
}

// The destination path on the remote host. The last element of the path will be created if it doesn't exist but it's an error when additional elements don't exist. When the remote path is an existing directory, the source file or directory will be copied into that directory. When the source is a file and the remote path is an existing file, that file will be overwritten. When the source is a directory and the remote path an existing file, the copy will fail.
func (o CopyToRemoteOutput) RemotePath() pulumi.StringOutput {
	return o.ApplyT(func(v *CopyToRemote) pulumi.StringOutput { return v.RemotePath }).(pulumi.StringOutput)
}

// An [asset or an archive](https://www.pulumi.com/docs/concepts/assets-archives/) to upload as the source of the copy. It must be path-based, i.e., be a `FileAsset` or a `FileArchive`. The item will be copied as-is; archives like .tgz will not be unpacked. Directories are copied recursively, overwriting existing files.
func (o CopyToRemoteOutput) Source() pulumi.AssetOrArchiveOutput {
	return o.ApplyT(func(v *CopyToRemote) pulumi.AssetOrArchiveOutput { return v.Source }).(pulumi.AssetOrArchiveOutput)
}

// Trigger replacements on changes to this input.
func (o CopyToRemoteOutput) Triggers() pulumi.ArrayOutput {
	return o.ApplyT(func(v *CopyToRemote) pulumi.ArrayOutput { return v.Triggers }).(pulumi.ArrayOutput)
}

type CopyToRemoteArrayOutput struct{ *pulumi.OutputState }

func (CopyToRemoteArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*CopyToRemote)(nil)).Elem()
}

func (o CopyToRemoteArrayOutput) ToCopyToRemoteArrayOutput() CopyToRemoteArrayOutput {
	return o
}

func (o CopyToRemoteArrayOutput) ToCopyToRemoteArrayOutputWithContext(ctx context.Context) CopyToRemoteArrayOutput {
	return o
}

func (o CopyToRemoteArrayOutput) Index(i pulumi.IntInput) CopyToRemoteOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *CopyToRemote {
		return vs[0].([]*CopyToRemote)[vs[1].(int)]
	}).(CopyToRemoteOutput)
}

type CopyToRemoteMapOutput struct{ *pulumi.OutputState }

func (CopyToRemoteMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*CopyToRemote)(nil)).Elem()
}

func (o CopyToRemoteMapOutput) ToCopyToRemoteMapOutput() CopyToRemoteMapOutput {
	return o
}

func (o CopyToRemoteMapOutput) ToCopyToRemoteMapOutputWithContext(ctx context.Context) CopyToRemoteMapOutput {
	return o
}

func (o CopyToRemoteMapOutput) MapIndex(k pulumi.StringInput) CopyToRemoteOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *CopyToRemote {
		return vs[0].(map[string]*CopyToRemote)[vs[1].(string)]
	}).(CopyToRemoteOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*CopyToRemoteInput)(nil)).Elem(), &CopyToRemote{})
	pulumi.RegisterInputType(reflect.TypeOf((*CopyToRemoteArrayInput)(nil)).Elem(), CopyToRemoteArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*CopyToRemoteMapInput)(nil)).Elem(), CopyToRemoteMap{})
	pulumi.RegisterOutputType(CopyToRemoteOutput{})
	pulumi.RegisterOutputType(CopyToRemoteArrayOutput{})
	pulumi.RegisterOutputType(CopyToRemoteMapOutput{})
}
