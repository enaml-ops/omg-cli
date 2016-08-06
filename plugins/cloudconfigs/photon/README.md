# Photon Cloud Config Plugin

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
photon-cloudconfigplugin-osx --az z1 \
--az z2 \
--az z3 \
--network-name-1 bosh \
--network-az-1 z1 \
--network-cidr-1 10.0.0.0/26 \
--network-gateway-1 10.0.0.1 \
--network-dns-1 169.254.169.254,8.8.8.8 \
--network-reserved-1 10.0.0.1-10.0.0.2,10.0.0.60-10.0.0.63 \
--network-static-1 10.0.0.4,10.0.0.10 \
--photon-network-name-1 test-vnet \
--network-name-2 concourse \
--network-az-2 z2 \
--network-cidr-2 10.0.0.64/26 \
--network-gateway-2 10.0.0.65 \
--network-dns-2 169.254.169.254,8.8.8.8 \
--network-reserved-2 10.0.0.65-10.0.0.70,10.0.0.122-10.0.0.127 \
--network-static-2 10.0.0.72,10.0.0.73,10.0.0.74,10.0.0.75 \
--photon-network-name-2 test-vnet
```
