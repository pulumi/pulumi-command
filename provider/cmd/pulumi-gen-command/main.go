// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	command "github.com/pulumi/pulumi-command/provider/pkg/provider"
	providerVersion "github.com/pulumi/pulumi-command/provider/pkg/version"
)

// copied from encoding/json for use with JSONMarshal above
func MarshalIndent(v any) ([]byte, error) {

	// json.Marshal normally escapes HTML. This one doesn't
	// https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	b := buffer.Bytes()

	// serialize and pretty print
	var buf bytes.Buffer
	prefix, indent := "", "    "
	err = json.Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	flag.Usage = func() {
		const usageFormat = "Usage: %s <schema-file>"
		fmt.Fprintf(flag.CommandLine.Output(), usageFormat, os.Args[0])
		flag.PrintDefaults()
	}

	var version string
	flag.StringVar(&version, "version", providerVersion.Version, "the provider version to record in the generated code")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		return
	}
	s, err := command.Schema(version)
	if err != nil {
		panic(err)
	}

	// sort keys
	var arg map[string]any
	err = json.Unmarshal([]byte(s), &arg)
	if err != nil {
		panic(err)
	}

	// remove version key
	delete(arg, "version")

	// use custom marshal indent to skip html escaping
	out, err := MarshalIndent(arg)
	if err != nil {
		panic(err)
	}

	schemaPath := args[0]
	err = os.WriteFile(schemaPath, out, 0600)
	if err != nil {
		panic(err)
	}
}
