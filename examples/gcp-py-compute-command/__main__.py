import base64
import pulumi
import pulumi_gcp as gcp
import pulumi_command as command

# Get the config ready to go.
config = pulumi.Config()

# If keyName is provided, an existing KeyPair is used, else if publicKey is provided a new KeyPair
# derived from the publicKey is created.
key_name = config.get('keyName')
public_key = config.get('publicKey')


# The privateKey associated with the selected key must be provided (either directly or base64 encoded),
# along with an optional passphrase if needed.
def decode_key(key):
    try:
        key = base64.b64decode(key.encode('ascii')).decode('ascii')
    except:
        pass

    if key.startswith('-----BEGIN RSA PRIVATE KEY-----'):
        return key

    return key.encode('ascii')


private_key = config.require_secret('privateKey').apply(decode_key)

svcacct = gcp.serviceaccount.Account("my-service-account",
                                     account_id="service-account",
                                     display_name="Service Account for Ansible")

svckey = gcp.serviceaccount.Key("my-service-key",
                                service_account_id=svcacct.name,
                                public_key_type="TYPE_X509_PEM_FILE")


addr = gcp.compute.address.Address('my-address')
compute_instance = gcp.compute.Instance(
    "my-instance",
    machine_type="f1-micro",
    boot_disk=gcp.compute.InstanceBootDiskArgs(
        initialize_params=gcp.compute.InstanceBootDiskInitializeParamsArgs(
            image="ubuntu-os-cloud/ubuntu-2004-lts"
        )
    ),
    network_interfaces=[gcp.compute.InstanceNetworkInterfaceArgs(
        network='default',
        access_configs=[gcp.compute.InstanceNetworkInterfaceAccessConfigArgs(
            nat_ip=addr.address
        )],
    )],
    service_account=gcp.compute.InstanceServiceAccountArgs(
        scopes=["https://www.googleapis.com/auth/cloud-platform"],
        email=svcacct.email
    ),
    metadata={
        'ssh-keys': f'user:{public_key}'
    },
)

conn = command.remote.ConnectionArgs(
    host=addr.address,
    private_key=private_key,
    user='user'
)

# Copy a config file to our server.
cp_config = command.remote.Copy(
    'config',
    connection=conn,
    local_asset='myapp.conf',
    remote_path='myapp.conf',
    opts=pulumi.ResourceOptions(depends_on=[compute_instance])
)

# Execute a basic command on our server.
cat_config = command.remote.Command(
    'cat-config',
    connection=conn,
    create='cat myapp.conf',
    opts=pulumi.ResourceOptions(depends_on=[cp_config])
)

# Export the server's IP and stdout from the command.
pulumi.export('publicIp', addr.address)
pulumi.export('catConfigStdout', cat_config.stdout)
