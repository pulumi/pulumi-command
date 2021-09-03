# command Pulumi Provider

Use cases:
* Post create hook
* Readiness/health check
* 

## Build and Test

```bash
# build and install the resource provider plugin
$ make build install

# test
$ cd examples/simple
$ yarn link @pulumi/command
$ yarn install
$ pulumi stack init test
$ pulumi up
```

