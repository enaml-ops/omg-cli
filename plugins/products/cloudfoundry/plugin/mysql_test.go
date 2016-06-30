package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-cli/plugins/products/cf-mysql/enaml-gen/mysql"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MySQL Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then it should return the error and exit", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
			})
			_, err := NewMySQLPartition(c)
			Ω(err).ShouldNot(BeNil())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var err error
		var mysqlPartition InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--mysql-ip", "1.0.10.1",
				"--mysql-ip", "1.0.10.2",
				"--mysql-network", "foundry-net",
				"--mysql-vm-type", "blah",
				"--mysql-disk-type", "blah",
				"--mysql-admin-password", "mysqladmin",
				"--mysql-bootstrap-username", "mysqlbootstrap",
				"--mysql-bootstrap-password", "mysqlbootstrappwd",
				"--syslog-address", "syslog-server",
				"--syslog-port", "10601",
				"--syslog-transport", "tcp",
			})
			mysqlPartition, err = NewMySQLPartition(c)
		})
		It("then it should not return an error", func() {
			Ω(err).Should(BeNil())
		})
		It("then it should allow the user to configure the mysql IPs", func() {
			ig := mysqlPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.10.1", "1.0.10.2"))
		})
		It("then it should have 2 instances", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := mysqlPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := mysqlPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 1 job", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(1))
		})
		It("then it should then have mysql job", func() {
			ig := mysqlPartition.ToInstanceGroup()
			job := ig.GetJobByName("mysql")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*mysql.Mysql)
			Ω(props.AdminPassword).Should(Equal("mysqladmin"))
			Ω(props.DatabaseStartupTimeout).Should(Equal(1200))
			Ω(props.MaxConnections).Should(Equal(1500))
			Ω(props.InnodbBufferPoolSize).Should(Equal(2147483648))
			Ω(props.BootstrapEndpoint.Username).Should(Equal("mysqlbootstrap"))
			Ω(props.BootstrapEndpoint.Password).Should(Equal("mysqlbootstrappwd"))
			Ω(props.SyslogAggregator.Address).Should(Equal("syslog-server"))
			Ω(props.SyslogAggregator.Port).Should(Equal(10601))
			Ω(props.SyslogAggregator.Transport).Should(Equal("tcp"))
			Ω(props.ClusterIps).Should(ConsistOf("1.0.10.1", "1.0.10.2"))
			Ω(props.SeededDatabases).Should(BeEmpty())
		})
		Context("when initialized WITH a complete set of arguments and seeded databases", func() {
			var err error
			var mysqlPartition InstanceGrouper
			BeforeEach(func() {
				plugin := new(Plugin)
				c := plugin.GetContext([]string{
					"cloudfoundry",
					"--stemcell-name", "cool-ubuntu-animal",
					"--az", "eastprod-1",
					"--mysql-ip", "1.0.10.1",
					"--mysql-ip", "1.0.10.2",
					"--mysql-network", "foundry-net",
					"--mysql-vm-type", "blah",
					"--mysql-disk-type", "blah",
					"--mysql-admin-password", "mysqladmin",
					"--mysql-bootstrap-username", "mysqlbootstrap",
					"--mysql-bootstrap-password", "mysqlbootstrappwd",
					"--syslog-address", "syslog-server",
					"--syslog-port", "10601",
					"--syslog-transport", "tcp",
					"--db-uaa-username", "uaa-username",
					"--db-uaa-password", "uaa-password",
					"--db-ccdb-username", "ccdb-username",
					"--db-ccdb-password", "ccdb-password",
					"--db-console-username", "console-username",
					"--db-console-password", "console-password",
				})
				mysqlPartition, err = NewMySQLPartition(c)
			})
			It("then it should not return an error", func() {
				Ω(err).Should(BeNil())
			})
			It("then it should then have mysql job with 3 seeded dbs", func() {
				ig := mysqlPartition.ToInstanceGroup()
				job := ig.GetJobByName("mysql")
				Ω(job).ShouldNot(BeNil())
				props, _ := job.Properties.(*mysql.Mysql)
				Ω(props.AdminPassword).Should(Equal("mysqladmin"))
				Ω(props.DatabaseStartupTimeout).Should(Equal(1200))
				Ω(props.MaxConnections).Should(Equal(1500))
				Ω(props.InnodbBufferPoolSize).Should(Equal(2147483648))
				Ω(props.BootstrapEndpoint.Username).Should(Equal("mysqlbootstrap"))
				Ω(props.BootstrapEndpoint.Password).Should(Equal("mysqlbootstrappwd"))
				Ω(props.SyslogAggregator.Address).Should(Equal("syslog-server"))
				Ω(props.SyslogAggregator.Port).Should(Equal(10601))
				Ω(props.SyslogAggregator.Transport).Should(Equal("tcp"))
				Ω(props.ClusterIps).Should(ConsistOf("1.0.10.1", "1.0.10.2"))
				Ω(len(props.SeededDatabases.([]MySQLSeededDatabase))).Should(Equal(3))
				for _, seededDB := range props.SeededDatabases.([]MySQLSeededDatabase) {
					switch seededDB.Name {
					case "uaa":
						Ω(seededDB.Username).Should(Equal("uaa-username"))
						Ω(seededDB.Password).Should(Equal("uaa-password"))
					case "ccdb":
						Ω(seededDB.Username).Should(Equal("ccdb-username"))
						Ω(seededDB.Password).Should(Equal("ccdb-password"))
					case "console":
						Ω(seededDB.Username).Should(Equal("console-username"))
						Ω(seededDB.Password).Should(Equal("console-password"))
					default:
						panic("Unexpected db")
					}
				}
			})
		})
	})
})
