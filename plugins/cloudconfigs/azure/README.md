# Azure Cloud Config Plugin

### this is actually a plugin which is vendored
### it will be compiled separately and bundled
### with the final release


## A simple example running from inside cfjump
```
omg deploy-cloudconfig \
  --bosh-url https://bosh.url.com \
  --bosh-port 25555 \
  --bosh-user director \
  --bosh-pass passwd_here \
  --ssl-ignore --print-manifest \
    azure-cloudconfigplugin-linux \
  --az z1 \
  --network-name-1 bosh \
  --network-az-1 z1 \
  --network-cidr-1 10.0.0.0/24 \
  --network-gateway-1 10.0.0.1 \
  --network-dns-1 168.63.129.16 \
  --network-reserved-1 10.0.0.1-10.0.0.9 \
  --network-static-1 10.0.0.10-10.0.0.20 \
  --azure-virtual-network-name-1 pcf-net \
  --azure-subnet-name-1 pcf

```
