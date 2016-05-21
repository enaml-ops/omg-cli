package azurecli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/deployments/bosh-init"
	"github.com/xchapter7x/lo"
)

func deployYaml(myYaml string, boshInitDeploy func(string)) {
	fmt.Println("deploying your bosh")
	content := []byte(myYaml)
	tmpfile, err := ioutil.TempFile("", "bosh-init-deployment")
	defer os.Remove(tmpfile.Name())

	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	boshInitDeploy(tmpfile.Name())
}

func checkRequired(name string, c *cli.Context) {
	if c.String(name) == "" {
		fmt.Println("Sorry you need to provide " + name)
		os.Exit(1)
	}
}

func GetFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "name", Value: "bosh", Usage: "the vm name to be created in your azure account"},
		cli.StringFlag{Name: "bosh-release-ver", Value: "256.2", Usage: "the version of the bosh release you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-private-ip", Value: "10.0.0.4", Usage: "the private ip for the bosh vm to be created in azure"},
		cli.StringFlag{Name: "bosh-cpi-release-ver", Value: "11", Usage: "the bosh cpi version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-ver", Value: "3169", Usage: "the go agent version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-release-sha", Value: "ff2f4e16e02f66b31c595196052a809100cfd5a8", Usage: "sha1 of the bosh release being used (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-cpi-release-sha", Value: "395fc05c11ead59711188ebd0a684842a03dc93d", Usage: "sha1 of the cpi release being used (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-sha", Value: "ff13c47ac7ce121dee6153c1564bd8965edf9f59", Usage: "sha1 of the go agent being use (found on bosh.io)"},
		cli.StringFlag{Name: "director-name", Value: "my-bosh", Usage: "the name of your director"},
		cli.StringFlag{Name: "azure-public-ip", Value: "", Usage: "the static/public ip of your bosh"},
		cli.StringFlag{Name: "azure-instance-size", Value: "Standard_D1", Usage: "the instance size of your bosh"},
		cli.StringFlag{Name: "azure-vnet", Value: "", Usage: "your azure vnet name"},
		cli.StringFlag{Name: "azure-subnet", Value: "", Usage: "your azure subnet name"},
		cli.StringFlag{Name: "azure-subscription-id", Value: "", Usage: "your azure subscription id"},
		cli.StringFlag{Name: "azure-tenant-id", Value: "", Usage: "azure tenant id"},
		cli.StringFlag{Name: "azure-client-id", Value: "", Usage: "azure client id"},
		cli.StringFlag{Name: "azure-client-secret", Value: "", Usage: "azure client secret"},
		cli.StringFlag{Name: "azure-resource-group", Value: "", Usage: "azure resource group"},
		cli.StringFlag{Name: "azure-storage-account", Value: "", Usage: "azure storage account"},
		cli.StringFlag{Name: "azure-security-group", Value: "", Usage: "azure security group"},
		cli.StringFlag{Name: "azure-ssh-pub-key-path", Value: "", Usage: "the path to your azure ssh public key"},
		cli.StringFlag{Name: "azure-ssh-user", Value: "", Usage: "azure ssh user"},
		cli.StringFlag{Name: "azure-environment", Value: "AzureCloud", Usage: "the name of your azure environment"},
		cli.StringFlag{Name: "azure-private-key-path", Value: "", Usage: "the path to your private bosh key"},
		cli.BoolFlag{Name: "print-manifest", Usage: "if you would simply like to output a manifest the set this flag as true."},
	}
}

func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		var publicKey string
		checkRequired("azure-public-ip", c)
		checkRequired("azure-vnet", c)
		checkRequired("azure-subnet", c)
		checkRequired("azure-subscription-id", c)
		checkRequired("azure-tenant-id", c)
		checkRequired("azure-client-id", c)
		checkRequired("azure-client-secret", c)
		checkRequired("azure-resource-group", c)
		checkRequired("azure-storage-account", c)
		checkRequired("azure-security-group", c)
		checkRequired("azure-ssh-pub-key-path", c)
		checkRequired("azure-ssh-user", c)
		checkRequired("azure-private-key-path", c)

		if keybytes, err := ioutil.ReadFile(c.String("azure-ssh-pub-key-path")); err != nil {
			lo.G.Error("error in reading pubkey file: ", c.String("azure-ssh-pub-key-path"), err)
			os.Exit(1)
		} else {
			publicKey = string(keybytes)
		}

		manifest := boshinit.NewAzureBosh(boshinit.BoshInitConfig{
			Name:                      c.String("name"),
			BoshReleaseVersion:        c.String("bosh-release-ver"),
			BoshPrivateIP:             c.String("bosh-private-ip"),
			BoshCPIReleaseVersion:     c.String("bosh-cpi-release-ver"),
			GoAgentVersion:            c.String("go-agent-ver"),
			BoshReleaseSHA:            c.String("bosh-release-sha"),
			BoshCPIReleaseSHA:         c.String("bosh-cpi-release-sha"),
			GoAgentSHA:                c.String("go-agent-sha"),
			BoshDirectorName:          c.String("director-name"),
			BoshInstanceSize:          c.String("azure-instance-size"),
			AzurePublicIP:             c.String("azure-public-ip"),
			AzureVnet:                 c.String("azure-vnet"),
			AzureSubnet:               c.String("azure-subnet"),
			AzureSubscriptionID:       c.String("azure-subscription-id"),
			AzureTenantID:             c.String("azure-tenant-id"),
			AzureClientID:             c.String("azure-client-id"),
			AzureClientSecret:         c.String("azure-client-secret"),
			AzureResourceGroup:        c.String("azure-resource-group"),
			AzureStorageAccount:       c.String("azure-storage-account"),
			AzureDefaultSecurityGroup: c.String("azure-security-group"),
			AzureSSHPubKey:            publicKey,
			AzureSSHUser:              c.String("azure-ssh-user"),
			AzureEnvironment:          c.String("azure-environment"),
			AzurePrivateKeyPath:       c.String("azure-private-key-path"),
		})

		if yamlString, err := enaml.Paint(manifest); err == nil {

			if c.Bool("print-manifest") {
				fmt.Println(yamlString)

			} else {
				deployYaml(yamlString, boshInitDeploy)
			}
		} else {
			e = err
		}
		return
	}
}
