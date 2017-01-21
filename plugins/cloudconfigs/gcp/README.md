# GCP Cloud Config Plugin

### this is actually a plugin which is vendored
### it will be compiled separately and bundled
### with the final release


## Example
```
./omg-osx deploy-cloudconfig --bosh-url https://bosh.url.com \
--bosh-port 25555 \
--bosh-user admin \
--bosh-pass admin \
--ssl-ignore \
--print-manifest \
gcp-cloudconfigplugin-osx --az z1 \
--az z2 \
--az z3 \
--gcp-availability-zone test1 \
--gcp-availability-zone test2 \
--gcp-availability-zone test3 \
--network-name-1 bosh \
--network-az-1 z1 \
--network-cidr-1 10.0.0.0/26 \
--network-gateway-1 10.0.0.1 \
--network-dns-1 169.254.169.254 \
--network-dns-1 8.8.8.8 \
--network-reserved-1 10.0.0.1-10.0.0.2 \
--network-reserved-1 10.0.0.60-10.0.0.63 \
--network-static-1 10.0.0.4 \
--network-static-1 10.0.0.10 \
--gcp-network-name-1 test-vnet \
--gcp-subnetwork-name-1 test-subnet-bosh-us-east1 \
--gcp-network-tag-1 nat-traverse,no-ip \
--network-name-2 concourse \
--network-az-2 z2 \
--network-cidr-2 10.0.0.64/26 \
--network-gateway-2 10.0.0.65 \
--network-dns-2 169.254.169.254 \
--network-dns-2 8.8.8.8 \
--network-reserved-2 10.0.0.65-10.0.0.70 \
--network-reserved-2 10.0.0.122-10.0.0.127 \
--network-static-2 10.0.0.72 \
--network-static-2 10.0.0.73 \
--network-static-2 10.0.0.74 \
--network-static-2 10.0.0.75 \
--gcp-network-name-2 test-vnet \
--gcp-subnetwork-name-2 test-subnet-concourse-us-east1-c \
--gcp-network-tag-2 nat-traverse \
--gcp-network-tag-2 no-ip
```
