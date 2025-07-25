# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins as _builtins
import warnings
import sys
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
if sys.version_info >= (3, 11):
    from typing import NotRequired, TypedDict, TypeAlias
else:
    from typing_extensions import NotRequired, TypedDict, TypeAlias
from .. import _utilities
from . import outputs
from ._enums import *
from ._inputs import *

__all__ = ['CommandArgs', 'Command']

@pulumi.input_type
class CommandArgs:
    def __init__(__self__, *,
                 connection: pulumi.Input['ConnectionArgs'],
                 add_previous_output_in_env: Optional[pulumi.Input[_builtins.bool]] = None,
                 create: Optional[pulumi.Input[_builtins.str]] = None,
                 delete: Optional[pulumi.Input[_builtins.str]] = None,
                 environment: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]] = None,
                 logging: Optional[pulumi.Input['Logging']] = None,
                 stdin: Optional[pulumi.Input[_builtins.str]] = None,
                 triggers: Optional[pulumi.Input[Sequence[Any]]] = None,
                 update: Optional[pulumi.Input[_builtins.str]] = None):
        """
        The set of arguments for constructing a Command resource.
        :param pulumi.Input['ConnectionArgs'] connection: The parameters with which to connect to the remote host.
        :param pulumi.Input[_builtins.bool] add_previous_output_in_env: If the previous command's stdout and stderr (as generated by the prior create/update) is
               injected into the environment of the next run as PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR.
               Defaults to true.
        :param pulumi.Input[_builtins.str] create: The command to run once on resource creation.
               
               If an `update` command isn't provided, then `create` will also be run when the resource's inputs are modified.
               
               Note that this command will not be executed if the resource has already been created and its inputs are unchanged.
               
               Use `local.runOutput` if you need to run a command on every execution of your program.
        :param pulumi.Input[_builtins.str] delete: The command to run on resource delettion.
               
               The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the stdout and stderr properties of the Command resource from previous create or update steps.
        :param pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]] environment: Additional environment variables available to the command's process.
               Note that this only works if the SSH server is configured to accept these variables via AcceptEnv.
               Alternatively, if a Bash-like shell runs the command on the remote host, you could prefix the command itself
               with the variables in the form 'VAR=value command'.
        :param pulumi.Input['Logging'] logging: If the command's stdout and stderr should be logged. This doesn't affect the capturing of
               stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
               outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.
        :param pulumi.Input[_builtins.str] stdin: Pass a string to the command's process as standard in
        :param pulumi.Input[Sequence[Any]] triggers: The resource will be updated (or replaced) if any of these values change.
               
               The trigger values can be of any type.
               
               If the `update` command was provided the resource will be updated, otherwise it will be replaced using the `create` command.
               
               Please see the resource documentation for examples.
        :param pulumi.Input[_builtins.str] update: The command to run when the resource is updated.
               
               If empty, the create command will be executed instead.
               
               Note that this command will not run if the resource's inputs are unchanged.
               
               Use `local.runOutput` if you need to run a command on every execution of your program.
               
               The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the `stdout` and `stderr` properties of the Command resource from previous create or update steps.
        """
        pulumi.set(__self__, "connection", connection)
        if add_previous_output_in_env is None:
            add_previous_output_in_env = True
        if add_previous_output_in_env is not None:
            pulumi.set(__self__, "add_previous_output_in_env", add_previous_output_in_env)
        if create is not None:
            pulumi.set(__self__, "create", create)
        if delete is not None:
            pulumi.set(__self__, "delete", delete)
        if environment is not None:
            pulumi.set(__self__, "environment", environment)
        if logging is not None:
            pulumi.set(__self__, "logging", logging)
        if stdin is not None:
            pulumi.set(__self__, "stdin", stdin)
        if triggers is not None:
            pulumi.set(__self__, "triggers", triggers)
        if update is not None:
            pulumi.set(__self__, "update", update)

    @_builtins.property
    @pulumi.getter
    def connection(self) -> pulumi.Input['ConnectionArgs']:
        """
        The parameters with which to connect to the remote host.
        """
        return pulumi.get(self, "connection")

    @connection.setter
    def connection(self, value: pulumi.Input['ConnectionArgs']):
        pulumi.set(self, "connection", value)

    @_builtins.property
    @pulumi.getter(name="addPreviousOutputInEnv")
    def add_previous_output_in_env(self) -> Optional[pulumi.Input[_builtins.bool]]:
        """
        If the previous command's stdout and stderr (as generated by the prior create/update) is
        injected into the environment of the next run as PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR.
        Defaults to true.
        """
        return pulumi.get(self, "add_previous_output_in_env")

    @add_previous_output_in_env.setter
    def add_previous_output_in_env(self, value: Optional[pulumi.Input[_builtins.bool]]):
        pulumi.set(self, "add_previous_output_in_env", value)

    @_builtins.property
    @pulumi.getter
    def create(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        The command to run once on resource creation.

        If an `update` command isn't provided, then `create` will also be run when the resource's inputs are modified.

        Note that this command will not be executed if the resource has already been created and its inputs are unchanged.

        Use `local.runOutput` if you need to run a command on every execution of your program.
        """
        return pulumi.get(self, "create")

    @create.setter
    def create(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "create", value)

    @_builtins.property
    @pulumi.getter
    def delete(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        The command to run on resource delettion.

        The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the stdout and stderr properties of the Command resource from previous create or update steps.
        """
        return pulumi.get(self, "delete")

    @delete.setter
    def delete(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "delete", value)

    @_builtins.property
    @pulumi.getter
    def environment(self) -> Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]]:
        """
        Additional environment variables available to the command's process.
        Note that this only works if the SSH server is configured to accept these variables via AcceptEnv.
        Alternatively, if a Bash-like shell runs the command on the remote host, you could prefix the command itself
        with the variables in the form 'VAR=value command'.
        """
        return pulumi.get(self, "environment")

    @environment.setter
    def environment(self, value: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]]):
        pulumi.set(self, "environment", value)

    @_builtins.property
    @pulumi.getter
    def logging(self) -> Optional[pulumi.Input['Logging']]:
        """
        If the command's stdout and stderr should be logged. This doesn't affect the capturing of
        stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
        outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.
        """
        return pulumi.get(self, "logging")

    @logging.setter
    def logging(self, value: Optional[pulumi.Input['Logging']]):
        pulumi.set(self, "logging", value)

    @_builtins.property
    @pulumi.getter
    def stdin(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        Pass a string to the command's process as standard in
        """
        return pulumi.get(self, "stdin")

    @stdin.setter
    def stdin(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "stdin", value)

    @_builtins.property
    @pulumi.getter
    def triggers(self) -> Optional[pulumi.Input[Sequence[Any]]]:
        """
        The resource will be updated (or replaced) if any of these values change.

        The trigger values can be of any type.

        If the `update` command was provided the resource will be updated, otherwise it will be replaced using the `create` command.

        Please see the resource documentation for examples.
        """
        return pulumi.get(self, "triggers")

    @triggers.setter
    def triggers(self, value: Optional[pulumi.Input[Sequence[Any]]]):
        pulumi.set(self, "triggers", value)

    @_builtins.property
    @pulumi.getter
    def update(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        The command to run when the resource is updated.

        If empty, the create command will be executed instead.

        Note that this command will not run if the resource's inputs are unchanged.

        Use `local.runOutput` if you need to run a command on every execution of your program.

        The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the `stdout` and `stderr` properties of the Command resource from previous create or update steps.
        """
        return pulumi.get(self, "update")

    @update.setter
    def update(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "update", value)


@pulumi.type_token("command:remote:Command")
class Command(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 add_previous_output_in_env: Optional[pulumi.Input[_builtins.bool]] = None,
                 connection: Optional[pulumi.Input[Union['ConnectionArgs', 'ConnectionArgsDict']]] = None,
                 create: Optional[pulumi.Input[_builtins.str]] = None,
                 delete: Optional[pulumi.Input[_builtins.str]] = None,
                 environment: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]] = None,
                 logging: Optional[pulumi.Input['Logging']] = None,
                 stdin: Optional[pulumi.Input[_builtins.str]] = None,
                 triggers: Optional[pulumi.Input[Sequence[Any]]] = None,
                 update: Optional[pulumi.Input[_builtins.str]] = None,
                 __props__=None):
        """
        A command to run on a remote host. The connection is established via ssh.

        ## Example Usage

        ### A Basic Example
        This program connects to a server and runs the `hostname` command. The output is then available via the `stdout` property.

        ```python
        import pulumi
        import pulumi_command as command

        config = pulumi.Config()
        server = config.require("server")
        user_name = config.require("userName")
        private_key = config.require("privateKey")
        hostname_cmd = command.remote.Command("hostnameCmd",
            create="hostname",
            connection=command.remote.ConnectionArgs(
                host=server,
                user=user_name,
                private_key=private_key,
            ))
        pulumi.export("hostname", hostname_cmd.stdout)
        ```

        ### Triggers
        This example defines several trigger values of various kinds. Changes to any of them will cause `cmd` to be re-run.

        ```python
        import pulumi
        import pulumi_command as command
        import pulumi_random as random

        foo = "foo"
        file_asset_var = pulumi.FileAsset("Pulumi.yaml")
        rand = random.RandomString("rand", length=5)
        local_file = command.local.Command("localFile",
            create="touch foo.txt",
            archive_paths=["*.txt"])

        cmd = command.remote.Command("cmd",
            connection=command.remote.ConnectionArgs(
                host="insert host here",
            ),
            create="echo create > op.txt",
            delete="echo delete >> op.txt",
            triggers=[
                foo,
                rand.result,
                file_asset_var,
                local_file.archive,
            ])
        ```

        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[_builtins.bool] add_previous_output_in_env: If the previous command's stdout and stderr (as generated by the prior create/update) is
               injected into the environment of the next run as PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR.
               Defaults to true.
        :param pulumi.Input[Union['ConnectionArgs', 'ConnectionArgsDict']] connection: The parameters with which to connect to the remote host.
        :param pulumi.Input[_builtins.str] create: The command to run once on resource creation.
               
               If an `update` command isn't provided, then `create` will also be run when the resource's inputs are modified.
               
               Note that this command will not be executed if the resource has already been created and its inputs are unchanged.
               
               Use `local.runOutput` if you need to run a command on every execution of your program.
        :param pulumi.Input[_builtins.str] delete: The command to run on resource delettion.
               
               The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the stdout and stderr properties of the Command resource from previous create or update steps.
        :param pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]] environment: Additional environment variables available to the command's process.
               Note that this only works if the SSH server is configured to accept these variables via AcceptEnv.
               Alternatively, if a Bash-like shell runs the command on the remote host, you could prefix the command itself
               with the variables in the form 'VAR=value command'.
        :param pulumi.Input['Logging'] logging: If the command's stdout and stderr should be logged. This doesn't affect the capturing of
               stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
               outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.
        :param pulumi.Input[_builtins.str] stdin: Pass a string to the command's process as standard in
        :param pulumi.Input[Sequence[Any]] triggers: The resource will be updated (or replaced) if any of these values change.
               
               The trigger values can be of any type.
               
               If the `update` command was provided the resource will be updated, otherwise it will be replaced using the `create` command.
               
               Please see the resource documentation for examples.
        :param pulumi.Input[_builtins.str] update: The command to run when the resource is updated.
               
               If empty, the create command will be executed instead.
               
               Note that this command will not run if the resource's inputs are unchanged.
               
               Use `local.runOutput` if you need to run a command on every execution of your program.
               
               The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the `stdout` and `stderr` properties of the Command resource from previous create or update steps.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: CommandArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        A command to run on a remote host. The connection is established via ssh.

        ## Example Usage

        ### A Basic Example
        This program connects to a server and runs the `hostname` command. The output is then available via the `stdout` property.

        ```python
        import pulumi
        import pulumi_command as command

        config = pulumi.Config()
        server = config.require("server")
        user_name = config.require("userName")
        private_key = config.require("privateKey")
        hostname_cmd = command.remote.Command("hostnameCmd",
            create="hostname",
            connection=command.remote.ConnectionArgs(
                host=server,
                user=user_name,
                private_key=private_key,
            ))
        pulumi.export("hostname", hostname_cmd.stdout)
        ```

        ### Triggers
        This example defines several trigger values of various kinds. Changes to any of them will cause `cmd` to be re-run.

        ```python
        import pulumi
        import pulumi_command as command
        import pulumi_random as random

        foo = "foo"
        file_asset_var = pulumi.FileAsset("Pulumi.yaml")
        rand = random.RandomString("rand", length=5)
        local_file = command.local.Command("localFile",
            create="touch foo.txt",
            archive_paths=["*.txt"])

        cmd = command.remote.Command("cmd",
            connection=command.remote.ConnectionArgs(
                host="insert host here",
            ),
            create="echo create > op.txt",
            delete="echo delete >> op.txt",
            triggers=[
                foo,
                rand.result,
                file_asset_var,
                local_file.archive,
            ])
        ```

        :param str resource_name: The name of the resource.
        :param CommandArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(CommandArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 add_previous_output_in_env: Optional[pulumi.Input[_builtins.bool]] = None,
                 connection: Optional[pulumi.Input[Union['ConnectionArgs', 'ConnectionArgsDict']]] = None,
                 create: Optional[pulumi.Input[_builtins.str]] = None,
                 delete: Optional[pulumi.Input[_builtins.str]] = None,
                 environment: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]] = None,
                 logging: Optional[pulumi.Input['Logging']] = None,
                 stdin: Optional[pulumi.Input[_builtins.str]] = None,
                 triggers: Optional[pulumi.Input[Sequence[Any]]] = None,
                 update: Optional[pulumi.Input[_builtins.str]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = CommandArgs.__new__(CommandArgs)

            if add_previous_output_in_env is None:
                add_previous_output_in_env = True
            __props__.__dict__["add_previous_output_in_env"] = add_previous_output_in_env
            if connection is None and not opts.urn:
                raise TypeError("Missing required property 'connection'")
            __props__.__dict__["connection"] = None if connection is None else pulumi.Output.secret(connection)
            __props__.__dict__["create"] = create
            __props__.__dict__["delete"] = delete
            __props__.__dict__["environment"] = environment
            __props__.__dict__["logging"] = logging
            __props__.__dict__["stdin"] = stdin
            __props__.__dict__["triggers"] = triggers
            __props__.__dict__["update"] = update
            __props__.__dict__["stderr"] = None
            __props__.__dict__["stdout"] = None
        secret_opts = pulumi.ResourceOptions(additional_secret_outputs=["connection"])
        opts = pulumi.ResourceOptions.merge(opts, secret_opts)
        replace_on_changes = pulumi.ResourceOptions(replace_on_changes=["triggers[*]"])
        opts = pulumi.ResourceOptions.merge(opts, replace_on_changes)
        super(Command, __self__).__init__(
            'command:remote:Command',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'Command':
        """
        Get an existing Command resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = CommandArgs.__new__(CommandArgs)

        __props__.__dict__["add_previous_output_in_env"] = None
        __props__.__dict__["connection"] = None
        __props__.__dict__["create"] = None
        __props__.__dict__["delete"] = None
        __props__.__dict__["environment"] = None
        __props__.__dict__["logging"] = None
        __props__.__dict__["stderr"] = None
        __props__.__dict__["stdin"] = None
        __props__.__dict__["stdout"] = None
        __props__.__dict__["triggers"] = None
        __props__.__dict__["update"] = None
        return Command(resource_name, opts=opts, __props__=__props__)

    @_builtins.property
    @pulumi.getter(name="addPreviousOutputInEnv")
    def add_previous_output_in_env(self) -> pulumi.Output[Optional[_builtins.bool]]:
        """
        If the previous command's stdout and stderr (as generated by the prior create/update) is
        injected into the environment of the next run as PULUMI_COMMAND_STDOUT and PULUMI_COMMAND_STDERR.
        Defaults to true.
        """
        return pulumi.get(self, "add_previous_output_in_env")

    @_builtins.property
    @pulumi.getter
    def connection(self) -> pulumi.Output['outputs.Connection']:
        """
        The parameters with which to connect to the remote host.
        """
        return pulumi.get(self, "connection")

    @_builtins.property
    @pulumi.getter
    def create(self) -> pulumi.Output[Optional[_builtins.str]]:
        """
        The command to run once on resource creation.

        If an `update` command isn't provided, then `create` will also be run when the resource's inputs are modified.

        Note that this command will not be executed if the resource has already been created and its inputs are unchanged.

        Use `local.runOutput` if you need to run a command on every execution of your program.
        """
        return pulumi.get(self, "create")

    @_builtins.property
    @pulumi.getter
    def delete(self) -> pulumi.Output[Optional[_builtins.str]]:
        """
        The command to run on resource delettion.

        The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the stdout and stderr properties of the Command resource from previous create or update steps.
        """
        return pulumi.get(self, "delete")

    @_builtins.property
    @pulumi.getter
    def environment(self) -> pulumi.Output[Optional[Mapping[str, _builtins.str]]]:
        """
        Additional environment variables available to the command's process.
        Note that this only works if the SSH server is configured to accept these variables via AcceptEnv.
        Alternatively, if a Bash-like shell runs the command on the remote host, you could prefix the command itself
        with the variables in the form 'VAR=value command'.
        """
        return pulumi.get(self, "environment")

    @_builtins.property
    @pulumi.getter
    def logging(self) -> pulumi.Output[Optional['Logging']]:
        """
        If the command's stdout and stderr should be logged. This doesn't affect the capturing of
        stdout and stderr as outputs. If there might be secrets in the output, you can disable logging here and mark the
        outputs as secret via 'additionalSecretOutputs'. Defaults to logging both stdout and stderr.
        """
        return pulumi.get(self, "logging")

    @_builtins.property
    @pulumi.getter
    def stderr(self) -> pulumi.Output[_builtins.str]:
        """
        The standard error of the command's process
        """
        return pulumi.get(self, "stderr")

    @_builtins.property
    @pulumi.getter
    def stdin(self) -> pulumi.Output[Optional[_builtins.str]]:
        """
        Pass a string to the command's process as standard in
        """
        return pulumi.get(self, "stdin")

    @_builtins.property
    @pulumi.getter
    def stdout(self) -> pulumi.Output[_builtins.str]:
        """
        The standard output of the command's process
        """
        return pulumi.get(self, "stdout")

    @_builtins.property
    @pulumi.getter
    def triggers(self) -> pulumi.Output[Optional[Sequence[Any]]]:
        """
        The resource will be updated (or replaced) if any of these values change.

        The trigger values can be of any type.

        If the `update` command was provided the resource will be updated, otherwise it will be replaced using the `create` command.

        Please see the resource documentation for examples.
        """
        return pulumi.get(self, "triggers")

    @_builtins.property
    @pulumi.getter
    def update(self) -> pulumi.Output[Optional[_builtins.str]]:
        """
        The command to run when the resource is updated.

        If empty, the create command will be executed instead.

        Note that this command will not run if the resource's inputs are unchanged.

        Use `local.runOutput` if you need to run a command on every execution of your program.

        The environment variables `PULUMI_COMMAND_STDOUT` and `PULUMI_COMMAND_STDERR` are set to the `stdout` and `stderr` properties of the Command resource from previous create or update steps.
        """
        return pulumi.get(self, "update")

