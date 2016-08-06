package cloudconfigs

import (
	"fmt"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
)

type BaseCloudConfig struct {
	Manifest *enaml.CloudConfigManifest
}

func GetDeploymentManifestBytes(provider CloudConfigProvider) ([]byte, error) {
	var manifest *enaml.CloudConfigManifest
	var err error
	var cloudConfigYml []byte
	if manifest, err = CreateCloudConfigManifest(provider); err != nil {
		return nil, err
	}
	if cloudConfigYml, err = yaml.Marshal(manifest); err != nil {
		return nil, err
	}
	return cloudConfigYml, nil
}

func CreateCloudConfigManifest(provider CloudConfigProvider) (*enaml.CloudConfigManifest, error) {
	var err error
	var azs []enaml.AZ
	var networks []enaml.DeploymentNetwork
	var vmTypes []enaml.VMType
	var diskTypes []enaml.DiskType
	var compilation *enaml.Compilation

	base := BaseCloudConfig{
		Manifest: &enaml.CloudConfigManifest{},
	}

	if azs, err = provider.CreateAZs(); err != nil {
		return nil, err
	}
	base.Manifest.AZs = azs

	if networks, err = provider.CreateNetworks(); err != nil {
		return nil, err
	}
	base.Manifest.Networks = networks

	if vmTypes, err = provider.CreateVMTypes(); err != nil {
		return nil, err
	}
	base.Manifest.VMTypes = vmTypes

	if diskTypes, err = provider.CreateDiskTypes(); err != nil {
		return nil, err
	}
	base.Manifest.DiskTypes = diskTypes

	if compilation, err = provider.CreateCompilation(); err != nil {
		return nil, err
	}
	base.Manifest.Compilation = compilation

	return base.Manifest, nil
}

func CreateFlagnameWithSuffix(name string, suffix int) (flagname string) {
	return name + "-" + strconv.Itoa(suffix)
}

func CheckRequiredLength(targetLength, index int, c *cli.Context, names ...string) error {
	var invalidNames []string
	for _, name := range names {
		formattedName := fmt.Sprintf(name, index)
		if len(c.StringSlice(formattedName)) != targetLength {
			invalidNames = append(invalidNames, formattedName)
		}
	}
	if len(invalidNames) > 0 {
		err := fmt.Errorf("Sorry you need to provide %s flags with %d element(s) to continue", invalidNames, targetLength)
		return err
	}
	return nil
}
