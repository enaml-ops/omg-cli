package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-cli/pluginlib/util"
	"github.com/xchapter7x/lo"
	"gopkg.in/yaml.v2"
)

func init() {
	RegisterInstanceGrouperFactory(NewConsulPartition)
	RegisterInstanceGrouperFactory(NewNatsPartition)
	RegisterInstanceGrouperFactory(NewEtcdPartition)
	//diego_database-partition
	RegisterInstanceGrouperFactory(NewNFSPartition)
	RegisterInstanceGrouperFactory(NewGoRouterPartition)
	RegisterInstanceGrouperFactory(NewMySQLProxyPartition)
	RegisterInstanceGrouperFactory(NewMySQLPartition)
	//cloud_controller-partition
	//ha_proxy-partition
	//clock_global-partition
	RegisterInstanceGrouperFactory(NewCloudControllerWorkerPartition)
	//uaa-partition
	RegisterInstanceGrouperFactory(NewDiegoBrainPartition)
	//diego_cell-partition
	//doppler-partition
	//loggregator_trafficcontroller-partition
}

//GetFlags -
func (s *Plugin) GetFlags() (flags []cli.Flag) {
	return []cli.Flag{
		// shared for all instance groups:
		cli.StringFlag{Name: "stemcell-name", Usage: "the name of your desired stemcell"},
		cli.StringSliceFlag{Name: "az", Usage: "list of AZ names to use"},
		cli.StringFlag{Name: "network", Usage: "the name of the network to use"},

		cli.StringSliceFlag{Name: "router-ip", Usage: "a list of the router ips you wish to use"},
		cli.StringFlag{Name: "router-vm-type", Usage: "the name of your desired vm size"},
		cli.StringFlag{Name: "router-ssl-cert", Usage: "the go router ssl cert, or a filename preceded by '@'"},
		cli.StringFlag{Name: "router-ssl-key", Usage: "the go router ssl key, or a filename preceded by '@'"},
		cli.StringFlag{Name: "router-user", Value: "router_status", Usage: "the username of the go-routers"},
		cli.StringFlag{Name: "router-pass", Usage: "the password of the go-routers"},
		cli.BoolFlag{Name: "router-enable-ssl", Usage: "enable or disable ssl on your routers"},
		cli.StringFlag{Name: "nats-user", Value: "nats", Usage: "username for your nats pool"},
		cli.StringFlag{Name: "nats-pass", Value: "nats-password", Usage: "password for your nats pool"},
		cli.StringSliceFlag{Name: "nats-machine-ip", Usage: "ip of a nats node vm"},
		cli.StringSliceFlag{Name: "etcd-machine-ip", Usage: "ip of a etcd node vm"},
		cli.StringFlag{Name: "metron-zone", Usage: "zone guid for the metron agent"},
		cli.StringFlag{Name: "metron-secret", Usage: "shared secret for the metron agent endpoint"},
		cli.StringSliceFlag{Name: "consul-ip", Usage: "a list of the consul ips you wish to use"},
		cli.StringFlag{Name: "consul-vm-type", Usage: "the name of your desired vm size for consul"},
		cli.StringSliceFlag{Name: "consul-encryption-key", Usage: "encryption key for consul"},
		cli.StringFlag{Name: "consul-ca-cert", Usage: "ca cert contents for consul"},
		cli.StringFlag{Name: "consul-agent-cert", Usage: "agent cert contents for consul"},
		cli.StringFlag{Name: "consul-agent-key", Usage: "agent key contents for consul"},
		cli.StringFlag{Name: "consul-server-cert", Usage: "server cert contents for consul"},
		cli.StringFlag{Name: "consul-server-key", Usage: "server key contents for consul"},
		cli.StringFlag{Name: "syslog-address", Usage: "address of syslog server"},
		cli.IntFlag{Name: "syslog-port", Usage: "port of syslog server"},
		cli.StringFlag{Name: "syslog-transport", Usage: "transport to syslog server"},
		cli.StringFlag{Name: "etcd-vm-type", Usage: "the name of your desired vm size for etcd"},
		cli.StringFlag{Name: "etcd-disk-type", Usage: "the name of your desired persistent disk type for etcd"},
		cli.StringFlag{Name: "nats-vm-type", Usage: "the name of your desired vm size for NATS"},
		cli.StringSliceFlag{Name: "nfs-ip", Usage: "a list of the nfs ips you wish to use"},
		cli.StringFlag{Name: "nfs-vm-type", Usage: "the name of your desired vm size for nfs"},
		cli.StringFlag{Name: "nfs-disk-type", Usage: "the name of your desired persistent disk type for nfs"},
		cli.StringSliceFlag{Name: "nfs-allow-from-network-cidr", Usage: "the network cidr you wish to allow connections to nfs from"},

		//Mysql Flags
		cli.StringSliceFlag{Name: "mysql-ip", Usage: "a list of the mysql ips you wish to use"},
		cli.StringFlag{Name: "mysql-vm-type", Usage: "the name of your desired vm size for mysql"},
		cli.StringFlag{Name: "mysql-disk-type", Usage: "the name of your desired persistent disk type for mysql"},
		cli.StringFlag{Name: "mysql-admin-password", Usage: "admin password for mysql"},
		cli.StringFlag{Name: "mysql-bootstrap-username", Usage: "bootstrap username for mysql"},
		cli.StringFlag{Name: "mysql-bootstrap-password", Usage: "bootstrap password for mysql"},

		//MySQL proxy flags
		cli.StringSliceFlag{Name: "mysql-proxy-ip", Usage: "a list of mysql proxy ips you wish to use"},
		cli.StringFlag{Name: "mysql-proxy-vm-type", Usage: "the name of your desired vm size for mysql proxy"},
		cli.StringFlag{Name: "mysql-proxy-external-host", Usage: "Host name of MySQL proxy"},
		cli.StringFlag{Name: "mysql-proxy-api-username", Usage: "Proxy API user name"},
		cli.StringFlag{Name: "mysql-proxy-api-password", Usage: "Proxy API password"},

		//CC Worker Partition Flags
		cli.StringFlag{Name: "cc-worker-vm-type", Usage: "the name of the desired vm type for cc worker"},
		cli.StringFlag{Name: "cc-staging-upload-user", Usage: "user name for staging upload"},
		cli.StringFlag{Name: "cc-staging-upload-password", Usage: "password for staging upload"},
		cli.StringFlag{Name: "cc-bulk-api-user", Usage: "user name for bulk api calls"},
		cli.StringFlag{Name: "cc-bulk-api-password", Usage: "password for bulk api calls"},
		cli.StringFlag{Name: "cc-db-encryption-key", Usage: "Cloud Controller DB encryption key"},
		cli.StringFlag{Name: "cc-internal-api-user", Usage: "user name for Internal API calls"},
		cli.StringFlag{Name: "cc-internal-api-password", Usage: "password for Internal API calls"},
		cli.StringFlag{Name: "system-domain", Usage: "System Domain"},
		cli.StringSliceFlag{Name: "app-domain", Usage: "Applications Domain"},
		cli.StringFlag{Name: "allow-app-ssh-access", Usage: "Allow SSH Access?"},
		cli.StringFlag{Name: "nfs-server-address", Usage: "NFS Server address"},
		cli.StringFlag{Name: "nfs-share-path", Usage: "NFS Share Path"},

		cli.StringFlag{Name: "db-uaa-username", Usage: "uaa db username"},
		cli.StringFlag{Name: "db-uaa-password", Usage: "uaa db password"},
		cli.StringFlag{Name: "db-ccdb-username", Usage: "ccdb db username"},
		cli.StringFlag{Name: "db-ccdb-password", Usage: "ccdb db password"},
		cli.StringFlag{Name: "db-console-username", Usage: "console db username"},
		cli.StringFlag{Name: "db-console-password", Usage: "console db password"},

		// Diego Brain
		cli.StringSliceFlag{Name: "diego-brain-ip", Usage: "a list of static IPs for the diego brain"},
		cli.StringFlag{Name: "diego-brain-vm-type", Usage: "the name of the desired vm type for the diego brain"},
		cli.StringFlag{Name: "diego-brain-disk-type", Usage: "the name of your desired persistent disk type for the diego brain"},
	}
}

//GetMeta -
func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "cloudfoundry",
	}
}

//GetProduct -
func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	c := pluginutil.NewContext(args, s.GetFlags())
	dm := enaml.NewDeploymentManifest([]byte(``))
	dm.SetName(DeploymentName)
	dm.AddRelease(enaml.Release{Name: CFReleaseName, Version: CFReleaseVersion})
	dm.AddStemcell(enaml.Stemcell{OS: StemcellName, Version: StemcellVersion, Alias: StemcellAlias})

	for _, factory := range factories {
		// create and validate all registered instance groups
		grouper := factory(c)
		if grouper.HasValidValues() {
			ig := grouper.ToInstanceGroup()
			lo.G.Debug("instance-group: ", ig)
			dm.AddInstanceGroup(ig)
		} else {
			b, _ := yaml.Marshal(grouper)
			lo.G.Panic("invalid values in instance group: ", string(b))
		}
	}

	return dm.Bytes()
}

//GetContext -
func (s *Plugin) GetContext(args []string) (c *cli.Context) {
	c = pluginutil.NewContext(args, s.GetFlags())
	return
}
