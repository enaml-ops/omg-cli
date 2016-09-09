package cloudconfigs

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

func CreateAZFlags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringSliceFlag{Name: "az", Usage: "az name"},
	}
	return flags
}

func CreateNetworkFlags(flags []cli.Flag, iaasNetworkFlagFunction func([]cli.Flag, int) []cli.Flag) []cli.Flag {
	for i := 1; i <= SupportedNetworkCount; i++ {
		flags = append(flags, &cli.StringFlag{Name: CreateFlagnameWithSuffix("network-name", i), Usage: "network name"})
		flags = append(flags, &cli.StringSliceFlag{Name: CreateFlagnameWithSuffix("network-az", i), Usage: fmt.Sprintf("az of network %d", i)})
		flags = append(flags, &cli.StringSliceFlag{Name: CreateFlagnameWithSuffix("network-cidr", i), Usage: fmt.Sprintf("range of network %d", i)})
		flags = append(flags, &cli.StringSliceFlag{Name: CreateFlagnameWithSuffix("network-gateway", i), Usage: fmt.Sprintf("gateway of network %d", i)})
		flags = append(flags, &cli.StringSliceFlag{Name: CreateFlagnameWithSuffix("network-dns", i), Usage: fmt.Sprintf("comma delimited list of DNS servers for network %d", i)})
		flags = append(flags, &cli.StringSliceFlag{Name: CreateFlagnameWithSuffix("network-reserved", i), Usage: fmt.Sprintf("comma delimited list of reserved network ranges for network %d", i)})
		flags = append(flags, &cli.StringSliceFlag{Name: CreateFlagnameWithSuffix("network-static", i), Usage: fmt.Sprintf("comma delimited list of static IP addresses for network %d", i)})
		flags = iaasNetworkFlagFunction(flags, i)
	}
	return flags
}
