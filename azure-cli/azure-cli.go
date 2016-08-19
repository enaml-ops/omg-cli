package azurecli

import (
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/xchapter7x/lo"
)

func GetFlags() []cli.Flag {
	boshdefaults := boshinit.GetAzureDefaults()

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
			lo.G.Error(e.Error())
			return e
		}
		var publicKey string
		if err := utils.CheckRequired(c, "azure-vnet", "azure-subnet", "azure-subscription-id", "azure-tenant-id",
			"azure-client-id", "azure-client-secret", "azure-resource-group",
			"azure-storage-account", "azure-security-group",
			"azure-ssh-pub-key-path",
			"azure-ssh-user",
			"azure-private-key-path"); err != nil {
			lo.G.Error(err.Error())
			return err
		}

		if keybytes, err := ioutil.ReadFile(c.String("azure-ssh-pub-key-path")); err != nil {
			lo.G.Error("error in reading pubkey file: ", c.String("azure-ssh-pub-key-path"), err)
			return err
		} else {
			publicKey = string(keybytes)
		}

		provider := boshinit.NewAzureIaaSProvider(boshinit.AzureInitConfig{
			AzureInstanceSize:         c.String("azure-instance-size"),
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

		if err := boshBase.HandleDeployment(provider, boshInitDeploy); err != nil {
			return err
		}

		return nil
	}
}
