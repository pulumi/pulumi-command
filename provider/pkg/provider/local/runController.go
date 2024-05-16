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

package local

import "context"

// This is the Call method. It takes a RunInputs parameter and runs the command specified in
// it.
func (*Run) Call(ctx context.Context, input RunInputs) (RunOutputs, error) {
	r := RunOutputs{RunInputs: input}
	err := run(ctx, input.Command, r.RunInputs.BaseInputs, &r.BaseOutputs, input.Logging)
	return r, err
}
