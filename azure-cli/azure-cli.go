package azurecli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
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
	boshdefaults := boshinit.BoshDefaults{
		CIDR:               "10.0.0.0/24",
		Gateway:            "10.0.0.1",
		DNS:                &cli.StringSlice{"168.63.129.16"},
		BoshReleaseVersion: "256.2",
		BoshReleaseSHA:     "ff2f4e16e02f66b31c595196052a809100cfd5a8",
		CPIReleaseVersion:  "11",
		CPIReleaseSHA:      "395fc05c11ead59711188ebd0a684842a03dc93d",
		GOAgentVersion:     "3262.4",
		GOAgentSHA:         "1ec76310cd99d4ad2dd2b239b3dfde09c609b292",
		PrivateIP:          "10.0.0.4",
		NtpServers:         &cli.StringSlice{"0.pool.ntp.org", "1.pool.ntp.org"},
		CPIName:            "cpi",
	}

	boshFlags := boshinit.BoshFlags(boshdefaults)
	azureFlags := []cli.Flag{
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
	}
	for _, flag := range azureFlags {
		boshFlags = append(boshFlags, flag)
	}
	return boshFlags
}

func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		var boshBase *boshinit.BoshBase
		if boshBase, e = boshinit.NewBoshBase(c); e != nil {
			return
		}
		var publicKey string
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
			BoshInstanceSize:          c.String("azure-instance-size"),
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
		}, boshBase)

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
