import pulumi_command as command

hello = command.local.Command(
    'foo',
    create='echo hello',
    logging=command.local.Logging.STDOUT
)
