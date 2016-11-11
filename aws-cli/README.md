Prior to running this must provision RDS, connect and issue the following commands:

```
create database registry;
create database uaa;
```

To Deploy Bosh director on AWS with MySQL RDS and S3 blobstore the following command is used.

```
./omg-osx aws --mode uaa \
--bosh-public-ip <elastic-ip> \
--use-external-db \
--database-driver mysql2 \
--database-scheme mysql \
--database-port 3306 \
--database-user <user for rd> \
--database-password <password for rds> \
--database-host <host name for rds> \
--use-external-blobstore \
--blobstore-bucket <your bucket name> \
--aws-subnet subnet-XXXXX \
--aws-pem-path <path-to>.pem \
--aws-keyname <keypair-name> \
--aws-access-key XXXXXXXXX \
--aws-secret XXXXXX \
--aws-security-group <security group per bosh.io>

```
