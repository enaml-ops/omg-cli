package boshinit

import (
	"github.com/enaml-ops/omg-cli/deployments/bosh-init/enaml-gen/aws_cpi"
	"github.com/enaml-ops/omg-cli/deployments/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/deployments/bosh-init/enaml-gen/postgres"
	"github.com/enaml-ops/omg-cli/deployments/bosh-init/enaml-gen/registry"
)

type BoshInitConfig struct {
	Name                              string
	BoshReleaseVersion                string
	BoshReleaseSHA                    string
	BoshPrivateIP                     string
	BoshCPIReleaseVersion             string
	BoshCPIReleaseSHA                 string
	GoAgentVersion                    string
	GoAgentSHA                        string
	BoshAvailabilityZone              string
	BoshInstanceSize                  string
	BoshDirectorName                  string
	AWSSubnet                         string
	AWSElasticIP                      string
	AWSPEMFilePath                    string
	AWSAccessKeyID                    string
	AWSSecretKey                      string
	AWSRegion                         string
	AWSSecurityGroups                 []string
	AzurePublicIP                     string
	AzureVnet                         string
	AzureSubnet                       string
	AzureSubscriptionID               string
	AzureTenantID                     string
	AzureClientID                     string
	AzureClientSecret                 string
	AzureResourceGroup                string
	AzureStorageAccount               string
	AzureDefaultSecurityGroup         string
	AzureSSHPubKey                    string
	AzureSSHUser                      string
	AzureEnvironment                  string
	AzurePrivateKeyPath               string
	VSphereAddress                    string
	VSphereUser                       string
	VSpherePassword                   string
	VSphereDatacenterName             string
	VSphereVMFolder                   string
	VSphereTemplateFolder             string
	VSphereDatastorePattern           string
	VSpherePersistentDatastorePattern string
	VSphereDiskPath                   string
	VSphereClusters                   []string
	VSphereNetworks                   []Network
}

type Network struct {
	Name    string
	Range   string
	Gateway string
	DNS     []string
}

type Rr registry.Registry
type Ar aws_cpi.Registry

type RegistryProperty struct {
	Rr      `yaml:",inline"`
	Ar      `yaml:",inline"`
	Address string `yaml:"address"`
}
type user struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type DirectorProperty struct {
	director.Director `yaml:",inline"`
	Address           string
}

type PgSql struct {
	User     string
	Host     string
	Password string
	Database string
	Adapter  string
}

type Postgres interface {
	GetDirectorDB() *director.Db
	GetRegistryDB() *registry.Db
	GetPostgresDB() postgres.Postgres
}
