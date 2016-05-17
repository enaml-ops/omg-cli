# bosh-install
single executable to install bosh on different targeted IaaS'


This simply composes bosh-init and enaml to create a simple bosh cli installer


## download here: https://github.com/bosh-ops/bosh-install/releases/latest


## install bosh on aws
*check the bosh docs to setup your vpc (https://bosh.io/docs/init-aws.html)*
```
bosh-install-osx aws --aws-subnet subnet-123456 --aws-elastic-ip 12.34.567.890 --aws-pem-path ~/boshstuff/bosh.pem --aws-access-key  xxxxxxxxxxxx --aws-secret xxxxxxxxxx --aws-instance-size t2.micro --aws-region us-east-1 --aws-availability-zone us-east-1c
```

## AWS available options
```
bosh-install-osx aws --help
NAME:
   bosh-install-osx aws - aws [--flags] - deploy a bosh to aws

USAGE:
   bosh-install-osx aws [command options] [arguments...]

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
