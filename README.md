# - omg -   
#### it's kind of like an (o)ps (m)anager in (g)olang
an iaas independent plugable executable to install bosh, cloud configs and product deployments

[![wercker status](https://app.wercker.com/status/429f96482fd95fecbc70ecc25aee8c70/s/master "wercker status")](https://app.wercker.com/project/bykey/429f96482fd95fecbc70ecc25aee8c70)

[![release info](https://img.shields.io/github/downloads/enaml-ops/omg-cli/total.svg?maxAge=2592000 "release info")](http://www.somsubhra.com/github-release-stats/?username=enaml-ops&repository=omg-cli)

### What is OMG (http://enaml.pezapp.io)
omg is a cli tool. It natively allows users to:
- spin up a bosh on a target iaas,
- load it up with a cloud config
- deploy 'products' via their new bosh (vault, cloudfoundry, concourse, etc)

### What are plugins
#### plugins are your way of extending omg, providing a deployment definition or cloud config definition. instead of dealing with yaml or tiles, we build testable plugins using `enaml` and simply register them with omg.
*download a bundled plugin from a omg release or build your own*
*available plugin types are `cloudconfig` or `product` for more info about how to build a plugin take a look at one of the bundled plugins (ie. https://github.com/enaml-ops/omg-cli/tree/master/cloudconfigs/aws)*

### how we do bosh / cloud config / deployments
composes bosh-init, enaml and plugins to create a simple cli installer

### downloads:
- omg/cloudconfig (http://enaml.pezapp.io/downloads/release/)
  - omg-cli (azure, gcp, vsphere, aws, photon)
- product plugins (http://enaml.pezapp.io/downloads/products/)
  - vault, concourse, pcf, redis, docker

## install a BOSH using OMG-cli (aws example)
*check the bosh docs to setup your vpc (https://bosh.io/docs/init-aws.html)*
```bash
# download omg cli
$> wget -O omg https://github.com/enaml-ops/omg-cli/releases/download/v0.0.25/omg-osx && chmod +x omg
```

```
# the below dependencies only apply if you are looking to install a BOSH Director with omg-cli
$> sudo apt-get update
$> sudo apt-get install -y build-essential zlibc zlib1g-dev ruby ruby-dev openssl libxslt-dev libxml2-dev libssl-dev libreadline6 libreadline6-dev libyaml-dev libsqlite3-dev sqlite3

or 

$> xcode-select --install
$> brew install openssl
```

```bash
# deploy your bosh using the omg cli
$> ./omg aws \
--mode uaa \
--aws-subnet subnet-xxxxxxxxxxx \
--bosh-public-ip x.x.x.x \
--aws-pem-path ~/bosh.pem \
--aws-access-key  xxxxxxxxxxxxxxxxxxxxxx \
--aws-secret xxxxxxxxxxxxxxxxxxx \
--aws-instance-size t2.micro \
--aws-region us-east-1 \
--aws-availability-zone us-east-1c
```

*instructions on how to install BOSH on other supported iaas can be found by:*
```bash
$> ./omg azure --help
$> ./omg aws --help
$> ./omg vsphere --help
$> ./omg photon --help
$> ./omg gcp --help
```

## Setup Cloud Config on your BOSH (aws example)
```bash
# download cloudconfig plugin for aws
$> wget https://github.com/enaml-ops/omg-cli/releases/download/v0.0.25/aws-cloudconfigplugin-osx
```
```bash
# register the cloud config plugin for your iaas
$> ./omg register-plugin \
--type cloudconfig \
--pluginpath aws-cloudconfigplugin-osx
```

```bash
# to see your newly added plugin
$> ./omg list-cloudconfigs
Cloud Configs:
aws  -  .plugins/cloudconfig/aws-cloudconfigplugin-osx  -  map[]
```

```bash
# upload cloud config
$> ./omg deploy-cloudconfig \
--bosh-url https://bosh.url.com --bosh-port 25555 \
--bosh-user admin --bosh-pass admin --ssl-ignore \
aws-cloudconfigplugin-osx \
--aws-region us-east-1 \
--aws-security-group bosh \
--bosh-az-name-1 z1 \
--cidr-1 10.0.0.0/24 \
--gateway-1 10.0.0.1 \
--dns-1 10.0.0.2 \
--aws-az-name-1 us-east-1a \
--aws-subnet-name-1 aws-subnet-1 \
--bosh-reserve-range-1 "10.0.0.1-10.0.0.10"

```

*for information on other options and flags:*
```bash
$> ./omg deploy-cloudconfig aws-cloudconfigplugin-osx --help
```

## How to use omg + plugins to install a product (ex,. concourse on aws)

*tips & tricks*
- set `LOG_LEVEL=debug` for verbose output
- adding the `--print-manifest` flag with the bosh creds will simply print the manifest you are about to deploy

### bosh deployed concourse
*deploy a concourse*
```bash
# download concourse product plugin
$> wget https://github.com/enaml-ops/omg-product-bundle/releases/download/v0.0.14/concourse-plugin-osx
```

```bash
# register concourse product plugin
$> ./omg register-plugin --type product --pluginpath concourse-plugin-osx
```

```bash
# please only upload your releases and stemcells manually if your deployment does not use remote urls
# otherwise this will be automatically uploaded via omg-cli
```

```bash
# deploy your concourse
$> ./omg deploy-product \
--bosh-url https://bosh.url.com --bosh-port 25555 --bosh-user admin \
--bosh-pass admin --ssl-ignore \
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
