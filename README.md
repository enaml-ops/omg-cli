# - omg -   
#### it's kind of like an (o)ps (m)anager in (g)olang
an iaas independent plugable executable to install bosh, cloud configs and product deployments

[![wercker status](https://app.wercker.com/status/c2ef4a65c6f9b1f4d6292529b3c6fd77/s/master "wercker status")](https://app.wercker.com/project/bykey/c2ef4a65c6f9b1f4d6292529b3c6fd77)

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
#### plugins are your way of extending omg, providing a deployment definition or cloud config definition. instead of dealing with yaml or tiles, we build testable plugins using `enaml` and simply register them with omg.
*download a bundled plugin from a omg release or build your own*
*available plugin types are `cloudconfig` or `product` for more info about how to build a plugin take a look at one of the bundled plugins (ie. https://github.com/enaml-ops/omg-cli/tree/master/cloudconfigs/aws)*
```
$ ./omg-osx register-plugin --type cloudconfig --pluginpath ~/Downloads/aws-cloudconfigplugin-osx

# to see your newly added plugin
$ ./omg-osx list-cloudconfigs
Cloud Configs:
aws  -  .plugins/cloudconfig/aws-cloudconfigplugin-osx  -  map[]
```

## How to use omg + plugins to install concourse (bosh, cloud-config, aws and osx)

*tips & tricks*
- set `LOG_LEVEL=debug` for verbose output
- adding the `--print-manifest` flag with the bosh creds will simply print the manifest you are about to deploy

### initial setup
*install your omg-cli & plugins*
```
export VERSION=v0.0.12
export OS=osx
$ wget https://github.com/enaml-ops/omg-cli/releases/download/${VERSION}/omg-${OS}
$ wget https://github.com/enaml-ops/omg-cli/releases/download/${VERSION}/concourse-plugin-${OS}
$ wget https://github.com/enaml-ops/omg-cli/releases/download/${VERSION}/aws-cloudconfigplugin-${OS}

$ mv ./omg-${OS} omg && chmod +x omg
$ ./omg register-plugin --type cloudconfig --pluginpath aws-cloudconfigplugin-${OS}
$ ./omg list-cloudconfigs
$ ./omg register-plugin --type product --pluginpath concourse-plugin-${OS}
$ ./omg list-products
```

### bosh install
*build your bosh*
```
$ ./omg aws \
--aws-subnet subnet-123456 \
--aws-elastic-ip bosh.url.com \
--aws-pem-path ~/boshstuff/bosh.pem \
--aws-access-key  xxxxxxxxxxxxxxxxxxxx \
--aws-secret xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
--aws-instance-size t2.micro \
--aws-region us-east-1 \
--aws-availability-zone us-east-1c
```

### setup cloud config
*setup a cloud config*
```
$ ./omg deploy-cloudconfig \
--bosh-url https://bosh.url.com \
--bosh-port 25555 \
--bosh-user admin \
--bosh-pass admin \
--ssl-ignore \
--print-manifest \
aws-cloudconfigplugin-osx \
--aws-region us-east-1 \
--aws-security-group bosh \
--bosh-az-name-1 z1 \
--aws-az-name-1 us-east-1a \
--cidr-1 10.10.0.0/24 \
--gateway-1 10.0.0.1 \
--aws-subnet-name-1 aws-subnet-1 \
--dns-1 10.10.0.2 \
--bosh-reserve-range-1 "10.10.0.3-10.10.0.25" \
--bosh-az-name-2 z2 \
--aws-az-name-2 us-east-1b \
--cidr-2 10.10.64.0/24 \
--gateway-2 10.10.64.1 \
--aws-subnet-name-2 aws-subnet-2 \
--dns-2 10.10.0.2 \
--bosh-reserve-range-2 "10.10.64.3-10.10.64.25"
```

### bosh deployed concourse
*deploy a concourse*
```
# please only upload your releases and stemcells manually if your deployment does not use remote urls
# otherwise this will be automatically uploaded via omg-cli
$ bosh upload release https://bosh.io/d/github.com/concourse/concourse?v=1.0.1
$ bosh upload release https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.337.0
$ bosh upload stemcell https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=3232.4

$ ./omg deploy-product \
--bosh-url https://bosh.url.com \
--bosh-port 25555 \
--bosh-user admin \
--bosh-pass admin \
--ssl-ignore \
concourse-plugin-osx \
--web-vm-type small \
--worker-vm-type small \
--database-vm-type small \
--network-name private \
--url my.concourse.com \
--username concourse \
--password concourse \
--web-instances 1 \
--web-azs us-east-1c \
--worker-azs us-east-1c \
--database-azs us-east-1c \
--bosh-stemcell-alias trusty \
--postgresql-db-pwd secret \
--database-storage-type medium
```
