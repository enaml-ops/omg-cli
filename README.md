# - omg -   
#### it's kind of like an (o)ps (m)anager in (g)olang
an iaas independent plugable executable to install bosh, cloud configs and product deployments

[![wercker status](https://app.wercker.com/status/3ebf8db4be00a8cc9fb51fc669ed6026/s/master "wercker status")](https://app.wercker.com/project/bykey/3ebf8db4be00a8cc9fb51fc669ed6026)

### how we do bosh / cloud config / deployments
composes bosh-init, enaml and plugins to create a simple cli installer


## download here: 
https://github.com/enaml-ops/omg-cli/releases/latest

## install bosh on aws
*check the bosh docs to setup your vpc (https://bosh.io/docs/init-aws.html)*
```
$ omg-osx aws --aws-subnet subnet-123456 --aws-elastic-ip 12.34.567.890 --aws-pem-path ~/boshstuff/bosh.pem --aws-access-key  xxxxxxxxxxxx --aws-secret xxxxxxxxxx --aws-instance-size t2.micro --aws-region us-east-1 --aws-availability-zone us-east-1c
```

## install bosh on azure
*check the bosh docs to setup your vpc (https://bosh.io/docs/init-azure.html)*
```
$ $ ./omg-osx azure --name bosh --azure-public-ip xxxx --azure-vnet xxxx --azure-subnet xxxx --azure-subscription-id xxxx --azure-tenant-id xxxx --azure-client-id xxxx --azure-client-secret xxxx --azure-resource-group xxxx --azure-storage-account xxxx --azure-security-group xxxx --azure-ssh-pub-key xxxx --azure-ssh-user xxxx --azure-private-key-path xxxx
```

## register a plugin
### plugins are your way of extending omg, providing a deployment definition or cloud config definition. instead of dealing with yaml or tiles, we build testable plugins and simply register them with omg.
*download a bundled plugin from a omg release or build your own*
*available plugin types are `cloudconfig` or `product` for more info about how to build a plugin take a look at one of the bundled plugins (ie. https://github.com/enaml-ops/omg-cli/tree/master/cloudconfigs/aws)*
```
$ ./omg-osx register-plugin --type cloudconfig --pluginpath ~/Downloads/aws-cloudconfigplugin-osx

# to see your newly added plugin
$ ./omg-osx list-cloudconfigs
Cloud Configs:
aws  -  .plugins/cloudconfig/aws-cloudconfigplugin-osx  -  map[]
```




## List available options using `--help` or `-h`
ie.
```
$ ./omg-osx aws --help
NAME:
   omg-osx aws - aws [--flags] - deploy a bosh to aws

USAGE:
   omg-osx aws [command options] [arguments...]

OPTIONS:
   --name value                the vm name to be created in your ec2 account (default: "bosh")
   --bosh-release-ver value        the version of the bosh release you wish to use (found on bosh.io) (default: "256.2")
   --bosh-private-ip value        the private ip for the bosh vm to be created in ec2 (default: "10.0.0.6")
   --bosh-cpi-release-ver value        the bosh cpi version you wish to use (found on bosh.io) (default: "52")
   --go-agent-ver value            the go agent version you wish to use (found on bosh.io) (default: "3012")
   --bosh-release-sha value        sha1 of the bosh release being used (found on bosh.io) (default: "ff2f4e16e02f66b31c595196052a809100cfd5a8")
   --bosh-cpi-release-sha value        sha1 of the cpi release being used (found on bosh.io) (default: "dc4a0cca3b33dce291e4fbeb9e9948b6a7be3324")
   --go-agent-sha value            sha1 of the go agent being use (found on bosh.io) (default: "3380b55948abe4c437dee97f67d2d8df4eec3fc1")
   --aws-instance-size value        the size of aws instance you wish to create (default: "m3.xlarge")
   --aws-availability-zone value    the ec2 az you wish to deploy to (default: "us-east-1c")
   --director-name value        the name of your director (default: "my-bosh")
   --aws-subnet value            your target vpc subnet
   --aws-elastic-ip value        your elastic ip to assign to the bosh vm
   --aws-pem-path value            your aws pem file path
   --aws-access-key value        aws account access key
   --aws-secret value            aws account secret key
   --aws-region value            ec2 region to deploy on (default: "us-east-1")
   --print-manifest            if you would simply like to output a manifest the set this flag as true.
   ```
