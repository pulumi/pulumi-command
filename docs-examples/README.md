# Docs HCL Examples

Look-aside copies of the HCL snippets that appear in the Pulumi Command
provider's registry docs (`docs/_index.md`). Each subdirectory contains a
self-contained Pulumi HCL program built around a single example from the
docs page.

These programs were used to validate the HCL syntax in the docs against the
real `pulumi-language-hcl` runtime and the Pulumi Command provider schema.
They are not part of the published examples directory and are not exercised
by CI.

| Directory | Example in `docs/_index.md` | Verification |
|-----------|-----------------------------|--------------|
| `random-local-command/` | "A simple local resource (random)" | `pulumi up` end-to-end (no cloud creds required) |
| `copy-to-remote/` | "Remote Commands and Copying Assets To Remote Hosts" | `pulumi preview` with mock SSH config |
| `invoke-lambda/` | "Invoking a Lambda during Pulumi deployment" | `pulumi preview` against the AWS provider with `skipCredentialsValidation` |
| `eks-cleanup/` | "Graceful cleanup of workloads in a Kubernetes cluster" | HCL parses and resolves `eks_cluster` and `command_local_command`; the EKS component itself calls real AWS APIs during construction so a full preview requires creds |

## Running an example

Each project sets `runtime: hcl` in `Pulumi.yaml`. To run one locally:

```bash
cd random-local-command
pulumi stack init dev
pulumi up
```

`pulumi-language-hcl` and `pulumi-converter-hcl` must be on `$PATH`. See
[pulumi-labs/pulumi-hcl](https://github.com/pulumi-labs/pulumi-hcl) for
installation instructions.
