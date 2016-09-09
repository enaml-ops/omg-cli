package azurecli

import (
	"io/ioutil"

	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

func GetFlags() []pcli.Flag {
	boshdefaults := boshinit.GetAzureDefaults()

	boshFlags := boshinit.BoshFlags(boshdefaults)
	azureFlags := []pcli.Flag{
		pcli.CreateStringFlag("azure-instance-size", "the instance size of your bosh", "Standard_D1"),
		pcli.CreateStringFlag("azure-vnet", "your azure vnet name"),
		pcli.CreateStringFlag("azure-subnet", "your azure subnet name"),
		pcli.CreateStringFlag("azure-subscription-id", "your azure subscription id"),
		pcli.CreateStringFlag("azure-tenant-id", "azure tenant id"),
		pcli.CreateStringFlag("azure-client-id", "azure client id"),
		pcli.CreateStringFlag("azure-client-secret", "azure client secret"),
		pcli.CreateStringFlag("azure-resource-group", "azure resource group"),
		pcli.CreateStringFlag("azure-storage-account", "azure storage account"),
		pcli.CreateStringFlag("azure-security-group", "azure security group"),
		pcli.CreateStringFlag("azure-ssh-pub-key-path", "the path to your azure ssh public key"),
		pcli.CreateStringFlag("azure-ssh-user", "azure ssh user"),
		pcli.CreateStringFlag("azure-environment", "the name of your azure environment", "AzureCloud"),
		pcli.CreateStringFlag("azure-private-key-path", "the path to your private bosh key"),
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
