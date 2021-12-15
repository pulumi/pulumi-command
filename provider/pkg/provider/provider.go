// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/mapper"

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	pbempty "github.com/golang/protobuf/ptypes/empty"
)

type commandProvider struct {
	host        *provider.HostClient
	name        string
	version     string
	cancelFuncs map[context.Context]context.CancelFunc
}

func makeProvider(host *provider.HostClient, name, version string) (pulumirpc.ResourceProviderServer, error) {
	// Return the new provider
	return &commandProvider{
		host:        host,
		name:        name,
		version:     version,
		cancelFuncs: make(map[context.Context]context.CancelFunc),
	}, nil
}

// Call dynamically executes a method in the provider associated with a component resource.
func (k *commandProvider) Call(ctx context.Context, req *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Call is not yet implemented")
}

// Construct creates a new component resource.
func (k *commandProvider) Construct(ctx context.Context, req *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Construct is not yet implemented")
}

// CheckConfig validates the configuration for this provider.
func (k *commandProvider) CheckConfig(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// DiffConfig diffs the configuration for this provider.
func (k *commandProvider) DiffConfig(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return &pulumirpc.DiffResponse{}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (k *commandProvider) Configure(_ context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	return &pulumirpc.ConfigureResponse{}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (k *commandProvider) Invoke(_ context.Context, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	tok := req.GetTok()
	return nil, fmt.Errorf("unknown Invoke token %q", tok)
}

// StreamInvoke dynamically executes a built-in function in the provider. The result is streamed
// back as a series of messages.
func (k *commandProvider) StreamInvoke(req *pulumirpc.InvokeRequest, server pulumirpc.ResourceProvider_StreamInvokeServer) error {
	tok := req.GetTok()
	return fmt.Errorf("unknown StreamInvoke token %q", tok)
}

func check(urn resource.URN) error {
	ty := urn.Type()
	valid := []string{"local:Command", "remote:Command", "remote:CopyFile"}
	for _, v := range valid {
		if "command:"+v == string(ty) {
			return nil
		}
	}
	return fmt.Errorf("unknown resource type %q", ty)
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (k *commandProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	if err := check(urn); err != nil {
		return nil, err
	}
	return &pulumirpc.CheckResponse{Inputs: req.News, Failures: nil}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (k *commandProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	if err := check(urn); err != nil {
		return nil, err
	}

	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	changes := pulumirpc.DiffResponse_DIFF_NONE
	var diffs, replaces []string
	properties := map[string]bool{
		"environment":      false,
		"dir":              false,
		"interpreter":      false,
		"create":           false,
		"delete":           false,
		"localPath":        false,
		"remotePath":       false,
		"connection":       true,
		"replaceOnChanges": true,
	}
	if d := olds.Diff(news); d != nil {
		for key, replace := range properties {
			i := sort.SearchStrings(req.IgnoreChanges, key)
			if i < len(req.IgnoreChanges) && req.IgnoreChanges[i] == key {
				continue
			}

			if d.Changed(resource.PropertyKey(key)) {
				changes = pulumirpc.DiffResponse_DIFF_SOME
				diffs = append(diffs, key)

				if replace {
					replaces = append(replaces, key)
				}
			}
		}
	}

	return &pulumirpc.DiffResponse{
		Changes:  changes,
		Diffs:    diffs,
		Replaces: replaces,
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
func (k *commandProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	ctx = k.addContext(ctx)
	defer k.removeContext(ctx)
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if err := check(urn); err != nil {
		return nil, err
	}

	inputProps, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	inputs := inputProps.Mappable()

	var id string
	var outputs map[string]interface{}

	switch ty {
	case "command:local:Command":
		var cmd command
		err = mapper.MapI(inputs, &cmd)
		if err != nil {
			return nil, err
		}

		id, err = cmd.RunCreate(ctx, k.host, urn)
		if err != nil {
			return nil, err
		}

		outputs, err = mapper.New(&mapper.Opts{IgnoreMissing: true, IgnoreUnrecognized: true}).Encode(cmd)
		if err != nil {
			return nil, err
		}
	case "command:remote:Command":
		var cmd remotecommand
		err = mapper.MapI(inputs, &cmd)
		if err != nil {
			return nil, err
		}

		id, err = cmd.RunCreate(ctx, k.host, urn)
		if err != nil {
			return nil, err
		}

		outputs, err = mapper.New(&mapper.Opts{IgnoreMissing: true, IgnoreUnrecognized: true}).Encode(cmd)
		if err != nil {
			return nil, err
		}
	case "command:remote:CopyFile":
		var cpf remotefilecopy
		err = mapper.MapI(inputs, &cpf)
		if err != nil {
			return nil, err
		}

		id, err = cpf.RunCreate(ctx, k.host, urn)
		if err != nil {
			return nil, err
		}

		outputs, err = mapper.New(&mapper.Opts{IgnoreMissing: true, IgnoreUnrecognized: true}).Encode(cpf)
		if err != nil {
			return nil, err
		}
	}

	outputProperties, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(outputs),
		plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}
	return &pulumirpc.CreateResponse{
		Id:         id,
		Properties: outputProperties,
	}, nil
}

// Read the current live state associated with a resource.
func (k *commandProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	ctx = k.addContext(ctx)
	defer k.removeContext(ctx)
	urn := resource.URN(req.GetUrn())
	if err := check(urn); err != nil {
		return nil, err
	}

	return &pulumirpc.ReadResponse{
		Id:         req.GetId(),
		Inputs:     req.GetInputs(),
		Properties: req.GetInputs(),
	}, nil
}

// Update updates an existing resource with new values.
func (k *commandProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	ctx = k.addContext(ctx)
	defer k.removeContext(ctx)
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if err := check(urn); err != nil {
		return nil, err
	}

	// Our Random resource will never be updated - if there is a diff, it will be a replacement.
	return nil, status.Errorf(codes.Unimplemented, "Update is not yet implemented for %q", ty)
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *commandProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	ctx = k.addContext(ctx)
	defer k.removeContext(ctx)
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if err := check(urn); err != nil {
		return nil, err
	}

	inputProps, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	inputs := inputProps.Mappable()

	decoder := mapper.New(&mapper.Opts{IgnoreMissing: true, IgnoreUnrecognized: true})

	switch ty {
	case "command:local:Command":
		var cmd command
		err = decoder.Decode(inputs, &cmd)
		if err != nil {
			return nil, err
		}

		err = cmd.RunDelete(ctx, k.host, urn)
		if err != nil {
			return nil, err
		}
	case "command:remote:Command":
		var cmd remotecommand
		err = decoder.Decode(inputs, &cmd)
		if err != nil {
			return nil, err
		}

		err = cmd.RunDelete(ctx, k.host, urn)
		if err != nil {
			return nil, err
		}
	case "command:remote:CopyFile":
		var cpf remotefilecopy
		err = decoder.Decode(inputs, &cpf)
		if err != nil {
			return nil, err
		}
		err = cpf.RunDelete(ctx, k.host, urn)
		if err != nil {
			return nil, err
		}
	}

	return &pbempty.Empty{}, nil
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (k *commandProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: k.version,
	}, nil
}

// GetSchema returns the JSON-serialized schema for the provider.
func (k *commandProvider) GetSchema(ctx context.Context, req *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	return &pulumirpc.GetSchemaResponse{}, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (k *commandProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	for _, cancel := range k.cancelFuncs {
		cancel()
	}
	return &pbempty.Empty{}, nil
}

func (k *commandProvider) addContext(c context.Context) context.Context {
	newctx, fn := context.WithCancel(c)
	k.cancelFuncs[newctx] = fn
	return newctx
}

func (k *commandProvider) removeContext(c context.Context) {
	delete(k.cancelFuncs, c)
}
