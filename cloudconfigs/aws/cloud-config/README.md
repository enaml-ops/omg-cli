
## How to work with cpis
```
$ enaml generate-jobs https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-aws-cpi-release?v=51
Could not find release in local cache. Downloading now.
27874562/27874562
completed generating release job structs for  https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-aws-cpi-release?v=51

$ enaml generate-jobs https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-openstack-cpi-release?v=24
Could not find release in local cache. Downloading now.
32781885/32781885
completed generating release job structs for  https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-openstack-cpi-release?v=24

$ enaml generate-jobs https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-vsphere-cpi-release?v=20
Could not find release in local cache. Downloading now.
29587401/29587401
completed generating release job structs for  https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-vsphere-cpi-release?v=20

$ ls enaml-gen/*
enaml-gen/aws_cpi:
agent.go             aws.go               awscpi.go            blobstore.go         connectionoptions.go env.go               nats.go              registry.go          stemcell.go

enaml-gen/openstack_cpi:
agent.go        blobstore.go    env.go          nats.go         openstack.go    openstackcpi.go registry.go

enaml-gen/vsphere_cpi:
agent.go             blobstore.go         connectionoptions.go env.go               nats.go              vcenter.go           vspherecpi.go
```
