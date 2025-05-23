// Copyright 2016-2022, Pulumi Corporation.
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

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	command "github.com/pulumi/pulumi-command/provider/pkg/provider"
	"github.com/pulumi/pulumi-command/provider/pkg/version"
)

// A provider is a program that listens for requests from the Pulumi engine
// to interact with cloud providers using a CRUD-based model.
func main() {
	version := strings.TrimPrefix(version.Version, "v")

	// This method defines the provider implemented in this repository.
	commandProvider := command.NewProvider()

	// This method starts serving requests using the Command provider.
	err := commandProvider.Run(context.Background(), "command", version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
