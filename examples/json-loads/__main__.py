import json
import pulumi
import pulumi_command as command

name = "foo"

vault_config = {
        "foo": "foo",
        "bar": "bar",
}

vault_ca = command.local.Command(
    "{}_vault_ca".format(name),
    create="cat",
    stdin=json.dumps(vault_config),
)
ca_secrets = vault_ca.stdout.apply(lambda x: json.loads(x))

pulumi.export("ca_secrets", ca_secrets)
