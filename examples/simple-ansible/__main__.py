import pulumi
import pulumi_command as command

local_command = command.local.Command(
    'ansible',
    create="ansible-playbook hello-world.yml"
)

pulumi.export('output', local_command.stdout)
