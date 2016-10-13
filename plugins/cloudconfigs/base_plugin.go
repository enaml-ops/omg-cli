package cloudconfigs

import (
	"fmt"

	"github.com/enaml-ops/pluginlib/pcli"
)

func CreateAZFlags() []pcli.Flag {
	return []pcli.Flag{
		pcli.CreateStringSliceFlag("az", "az name"),
	}
}

func CreateNetworkFlags(flags []pcli.Flag, iaasNetworkFlagFunction func([]pcli.Flag, int) []pcli.Flag) []pcli.Flag {
	for i := 1; i <= SupportedNetworkCount; i++ {
		flags = append(flags,
			pcli.CreateStringFlag(CreateFlagnameWithSuffix("network-name", i), "network name"),
			pcli.CreateStringSliceFlag(CreateFlagnameWithSuffix("network-az", i), fmt.Sprintf("az of network %d", i)),
			pcli.CreateStringSliceFlag(CreateFlagnameWithSuffix("network-cidr", i), fmt.Sprintf("range of network %d", i)),
			pcli.CreateStringSliceFlag(CreateFlagnameWithSuffix("network-gateway", i), fmt.Sprintf("gateway of network %d", i)),
			pcli.CreateStringSliceFlag(CreateFlagnameWithSuffix("network-dns", i), fmt.Sprintf("comma delimited list of DNS servers for network %d", i)),
			pcli.CreateStringSliceFlag(CreateFlagnameWithSuffix("network-reserved", i), fmt.Sprintf("comma delimited list of reserved network ranges for network %d", i)),
			pcli.CreateStringSliceFlag(CreateFlagnameWithSuffix("network-static", i), fmt.Sprintf("comma delimited list of static IP addresses for network %d", i)))

		flags = iaasNetworkFlagFunction(flags, i)
	}
	flags = append(flags, pcli.CreateBoolFlag("multi-assign-az", "Assigns all the AZs for each subnet"))
	return flags
}
