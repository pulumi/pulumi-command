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

package remote

import (
	"github.com/pulumi/pulumi-go-provider/infer"
)

type ProxyConnection struct{ connectionBase }

func (c *ProxyConnection) Annotate(a infer.Annotator) {
	a.Describe(&c, "Instructions for how to connect to a remote endpoint via a bastion host.")
	a.Describe(&c.User, "The user that we should use for the connection to the bastion host.")
	a.SetDefault(&c.User, "root")
	a.Describe(&c.Password, "The password we should use for the connection to the bastion host.")
	a.Describe(&c.Host, "The address of the bastion host to connect to.")
	a.Describe(&c.Port, "The port of the bastion host to connect to.")
	a.SetDefault(&c.Port, 22)
	a.Describe(&c.PrivateKey, "The contents of an SSH key to use for the connection. This takes preference over the password if provided.")
	a.Describe(&c.PrivateKeyPassword, "The password to use in case the private key is encrypted.")
	a.Describe(&c.AgentSocketPath, "SSH Agent socket path. Default to environment variable SSH_AUTH_SOCK if present.")
	a.Describe(&c.DialErrorLimit, "Max allowed errors on trying to dial the remote host. -1 set count to unlimited. Default value is 10.")
	a.SetDefault(&c.DialErrorLimit, dialErrorDefault)
	a.Describe(&c.PerDialTimeout, "Max number of seconds for each dial attempt. 0 implies no maximum. Default value is 15 seconds.")
	a.SetDefault(&c.PerDialTimeout, 15)
}
