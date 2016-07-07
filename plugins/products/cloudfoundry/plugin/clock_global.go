package cloudfoundry

import (
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/cloud_controller_clock"
	"github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/enaml-gen/cloud_controller_ng"
)

func NewClockGlobalPartition(c *cli.Context) InstanceGrouper {
	var db string
	mysqlProxies := c.StringSlice("mysql-proxy-ip")
	if len(mysqlProxies) > 0 {
		db = mysqlProxies[0]
	}

	cg := &clockGlobal{
		Instances:                1,
		AZs:                      c.StringSlice("az"),
		StemcellName:             c.String("stemcell-name"),
		VMTypeName:               c.String("clock-global-vm-type"),
		NetworkName:              c.String("network"),
		SystemDomain:             c.String("system-domain"),
		AppDomains:               c.StringSlice("app-domain"),
		Metron:                   NewMetron(c),
		Statsd:                   NewStatsdInjector(c),
		NFS:                      NewNFSMounter(c),
		AllowSSHAccess:           c.Bool("allow-app-ssh-access"),
		SkipSSLCertVerify:        c.BoolT("skip-cert-verify"),
		NATSUser:                 c.String("nats-user"),
		NATSPassword:             c.String("nats-pass"),
		NATSPort:                 c.Int("nats-port"),
		NATSMachines:             c.StringSlice("nats-machine-ip"),
		CloudController:          NewCloudControllerPartition(c).(*CloudControllerPartition),
		CCDBAddress:              db,
		JWTVerificationKey:       c.String("uaa-jwt-verification-key"),
		CCServiceDashboardSecret: c.String("cc-service-dashboards-client-secret"),
	}

	mysql := NewMySQLPartition(c).(*MySQL)
	ccdb := mysql.GetSeededDBByName("ccdb")
	if ccdb != nil {
		cg.CCDBUser = ccdb.Username
		cg.CCDBPassword = ccdb.Password
	}

	return cg
}

func (c *clockGlobal) ToInstanceGroup() *enaml.InstanceGroup {
	ig := &enaml.InstanceGroup{
		Name:      "clock_global-partition",
		Instances: 1,
		VMType:    c.VMTypeName,
		AZs:       c.AZs,
		Stemcell:  c.StemcellName,
		Networks: []enaml.Network{
			{Name: c.NetworkName},
		},
	}

	metronJob := c.Metron.CreateJob()
	nfsJob := c.NFS.CreateJob()
	statsdJob := c.Statsd.CreateJob()

	ccw := newCloudControllerNgWorkerJob(c.CloudController)
	props := ccw.Properties.(*cloud_controller_ng.CloudControllerNg)

	ig.AddJob(c.newCloudControllerClockJob(props))
	ig.AddJob(&metronJob)
	ig.AddJob(&nfsJob)
	ig.AddJob(&statsdJob)
	return ig
}

func (c *clockGlobal) newCloudControllerClockJob(ccng *cloud_controller_ng.CloudControllerNg) *enaml.InstanceJob {
	roles := make(map[string]string)
	roles["tag"] = "admin"
	roles["name"] = c.CCDBUser
	roles["password"] = c.CCDBPassword

	dbs := make(map[string]string)
	dbs["tag"] = "cc"
	dbs["name"] = "ccdb"
	dbs["citext"] = "true"

	props := &cloud_controller_clock.CloudControllerClock{
		Domain:                   c.SystemDomain,
		SystemDomain:             c.SystemDomain,
		SystemDomainOrganization: "system",
		AppDomains:               c.AppDomains,
		Cc:                       &cloud_controller_clock.Cc{},
		Ccdb: &cloud_controller_clock.Ccdb{
			Address:   c.CCDBAddress,
			Port:      3306,
			DbScheme:  "mysql",
			Roles:     roles,
			Databases: dbs,
		},
		Uaa: &cloud_controller_clock.Uaa{
			Url: prefixSystemDomain(c.SystemDomain, "uaa"),
			Jwt: &cloud_controller_clock.Jwt{
				VerificationKey: c.JWTVerificationKey,
			},
			Clients: &cloud_controller_clock.Clients{
				CcServiceDashboards: &cloud_controller_clock.CcServiceDashboards{
					Secret: c.CCServiceDashboardSecret,
				},
			},
		},
		LoggerEndpoint: &cloud_controller_clock.LoggerEndpoint{
			Port: strconv.Itoa(443),
		},
		Ssl: &cloud_controller_clock.Ssl{
			SkipCertVerify: c.SkipSSLCertVerify,
		},
		Nats: &cloud_controller_clock.Nats{
			User:     c.NATSUser,
			Password: c.NATSPassword,
			Port:     c.NATSPort,
			Machines: c.NATSMachines,
		},
	}

	job := &enaml.InstanceJob{
		Name:       "cloud_controller_clock",
		Release:    CFReleaseName,
		Properties: props,
	}

	ccYaml, _ := yaml.Marshal(ccng.Cc)
	yaml.Unmarshal(ccYaml, props.Cc)

	props.Cc.QuotaDefinitions = map[string]interface{}{
		"default": map[string]interface{}{
			"memory_limit":               10240,
			"total_services":             100,
			"non_basic_services_allowed": true,
			"total_routes":               1000,
			"trial_db_allowed":           true,
		},
		"runaway": map[string]interface{}{
			"memory_limit":               102400,
			"total_services":             -1,
			"non_basic_services_allowed": true,
			"total_routes":               1000,
		},
	}
	props.Cc.SecurityGroupDefinitions = []map[string]interface{}{
		map[string]interface{}{"name": "all_open",
			"rules": []map[string]interface{}{
				map[string]interface{}{
					"protocol":    "all",
					"destination": "0.0.0.0-255.255.255.255",
				},
			},
		},
	}
	return job
}

func (c *clockGlobal) HasValidValues() bool {
	return len(c.AZs) > 0 &&
		c.StemcellName != "" &&
		c.VMTypeName != "" &&
		c.NetworkName != "" &&
		c.Metron.HasValidValues() &&
		c.Statsd.HasValidValues() &&
		c.NFS.hasValidValues() &&
		c.CloudController.HasValidValues()
}
