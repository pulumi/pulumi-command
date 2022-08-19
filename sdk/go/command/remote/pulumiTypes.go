// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package remote

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Instructions for how to connect to a remote endpoint.
type Connection struct {
	// The address of the resource to connect to.
	Host string `pulumi:"host"`
	// The password we should use for the connection.
	Password *string `pulumi:"password"`
	// The port to connect to.
	Port *float64 `pulumi:"port"`
	// The contents of an SSH key to use for the connection. This takes preference over the password if provided.
	PrivateKey *string `pulumi:"privateKey"`
	// The user that we should use for the connection.
	User *string `pulumi:"user"`
}

// Defaults sets the appropriate defaults for Connection
func (val *Connection) Defaults() *Connection {
	if val == nil {
		return nil
	}
	tmp := *val
	if isZero(tmp.Port) {
		port_ := 22.0
		tmp.Port = &port_
	}
	if isZero(tmp.User) {
		user_ := "root"
		tmp.User = &user_
	}
	return &tmp
}

// ConnectionInput is an input type that accepts ConnectionArgs and ConnectionOutput values.
// You can construct a concrete instance of `ConnectionInput` via:
//
//	ConnectionArgs{...}
type ConnectionInput interface {
	pulumi.Input

	ToConnectionOutput() ConnectionOutput
	ToConnectionOutputWithContext(context.Context) ConnectionOutput
}

// Instructions for how to connect to a remote endpoint.
type ConnectionArgs struct {
	// The address of the resource to connect to.
	Host pulumi.StringInput `pulumi:"host"`
	// The password we should use for the connection.
	Password pulumi.StringPtrInput `pulumi:"password"`
	// The port to connect to.
	Port pulumi.Float64PtrInput `pulumi:"port"`
	// The contents of an SSH key to use for the connection. This takes preference over the password if provided.
	PrivateKey pulumi.StringPtrInput `pulumi:"privateKey"`
	// The user that we should use for the connection.
	User pulumi.StringPtrInput `pulumi:"user"`
}

// Defaults sets the appropriate defaults for ConnectionArgs
func (val *ConnectionArgs) Defaults() *ConnectionArgs {
	if val == nil {
		return nil
	}
	tmp := *val
	if isZero(tmp.Port) {
		tmp.Port = pulumi.Float64Ptr(22.0)
	}
	if isZero(tmp.User) {
		tmp.User = pulumi.StringPtr("root")
	}
	return &tmp
}
func (ConnectionArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*Connection)(nil)).Elem()
}

func (i ConnectionArgs) ToConnectionOutput() ConnectionOutput {
	return i.ToConnectionOutputWithContext(context.Background())
}

func (i ConnectionArgs) ToConnectionOutputWithContext(ctx context.Context) ConnectionOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ConnectionOutput)
}

// Instructions for how to connect to a remote endpoint.
type ConnectionOutput struct{ *pulumi.OutputState }

func (ConnectionOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Connection)(nil)).Elem()
}

func (o ConnectionOutput) ToConnectionOutput() ConnectionOutput {
	return o
}

func (o ConnectionOutput) ToConnectionOutputWithContext(ctx context.Context) ConnectionOutput {
	return o
}

// The address of the resource to connect to.
func (o ConnectionOutput) Host() pulumi.StringOutput {
	return o.ApplyT(func(v Connection) string { return v.Host }).(pulumi.StringOutput)
}

// The password we should use for the connection.
func (o ConnectionOutput) Password() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Connection) *string { return v.Password }).(pulumi.StringPtrOutput)
}

// The port to connect to.
func (o ConnectionOutput) Port() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v Connection) *float64 { return v.Port }).(pulumi.Float64PtrOutput)
}

// The contents of an SSH key to use for the connection. This takes preference over the password if provided.
func (o ConnectionOutput) PrivateKey() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Connection) *string { return v.PrivateKey }).(pulumi.StringPtrOutput)
}

// The user that we should use for the connection.
func (o ConnectionOutput) User() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Connection) *string { return v.User }).(pulumi.StringPtrOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ConnectionInput)(nil)).Elem(), ConnectionArgs{})
	pulumi.RegisterOutputType(ConnectionOutput{})
}
