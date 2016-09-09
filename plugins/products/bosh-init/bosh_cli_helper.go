package boshinit

import (
	"strconv"

	"github.com/enaml-ops/omg-cli/utils"
	"github.com/enaml-ops/pluginlib/pcli"
	"gopkg.in/urfave/cli.v2"
)

func BoshFlags(defaults *BoshBase) []pcli.Flag {
	return []pcli.Flag{
		pcli.CreateStringFlag("mode", "what type of bosh director to install.  Options are basic or uaa", "basic"),
		pcli.CreateStringFlag("cidr", "the network cidr range for your bosh deployment", defaults.NetworkCIDR),
		pcli.CreateStringFlag("gateway", "the gateway ip", defaults.NetworkGateway),
		pcli.CreateStringSliceFlag("dns", "the dns ip", defaults.NetworkDNS...),
		pcli.CreateStringFlag("bosh-private-ip", "the private ip for the bosh vm to be created", defaults.PrivateIP),
		pcli.CreateStringFlag("bosh-public-ip", "the public ip for the bosh vm to be created"),
		pcli.CreateStringFlag("bosh-release-sha", "sha1 of the bosh release being used (found on bosh.io)", defaults.BoshReleaseSHA),
		pcli.CreateStringFlag("bosh-release-url", "url to bosh release", defaults.BoshReleaseURL),
		pcli.CreateStringFlag("bosh-cpi-release-sha", "sha1 of the cpi release being used (found on bosh.io)", defaults.CPIReleaseSHA),
		pcli.CreateStringFlag("bosh-cpi-release-url", "url to bosh cpi release", defaults.CPIReleaseURL),
		pcli.CreateStringFlag("go-agent-release-sha", "sha1 of the go agent being use (found on bosh.io)", defaults.GOAgentSHA),
		pcli.CreateStringFlag("go-agent-release-url", "url to stemcell release", defaults.GOAgentReleaseURL),
		pcli.CreateStringFlag("director-name", "the name of your director", "enaml-bosh"),
		pcli.CreateStringFlag("uaa-release-sha", "sha1 of the uaa release being used (found on bosh.io)", "899f1e10f27e82ac524f1158a513392bbfabf2a0"),
		pcli.CreateStringFlag("uaa-release-url", "url to uaa release", "https://bosh.io/d/github.com/cloudfoundry/uaa-release?v=12.2"),
		pcli.CreateStringSliceFlag("ntp-server", "ntp server address", defaults.NtpServers...),
		pcli.CreateStringFlag("trusted-certs", "trusted ssl certs"),
		pcli.CreateStringFlag("nats-pwd", "password for nats"),
		pcli.CreateIntFlag("persistent-disk-size", "size of persistent disk", strconv.Itoa(defaults.PersistentDiskSize)),
		pcli.CreateBoolFlag("print-manifest", "if you would simply like to output a manifest the set this flag as true."),
		pcli.CreateStringFlag("hm-graphite-address", "graphite address to forward health monitor heartbeats"),
		pcli.CreateIntFlag("hm-graphite-port", "graphite port to forward health monitor heartbeats", "2003"),
		pcli.CreateStringFlag("syslog-address", "address of syslog server for forwarding heartbeats"),
		pcli.CreateIntFlag("syslog-port", "port of syslog server", "5514"),
		pcli.CreateStringFlag("syslog-transport", "transport to syslog server", "tcp"),
	}
}

var RequiredBoshFlags = []string{
	"cidr",
	"gateway",
	"dns",
	"bosh-private-ip",
	"bosh-release-url",
	"bosh-release-sha",
	"bosh-cpi-release-url",
	"bosh-cpi-release-sha",
	"go-agent-release-url",
	"go-agent-release-sha",
	"director-name",
	"uaa-release-url",
	"uaa-release-sha",
	"ntp-server",
	"persistent-disk-size",
}

func NewBoshBase(c *cli.Context) (base *BoshBase, err error) {

	utils.CheckRequired(c, RequiredBoshFlags...)

	base = &BoshBase{
		Mode:               c.String("mode"),
		NetworkCIDR:        c.String("cidr"),
		NetworkGateway:     c.String("gateway"),
		NetworkDNS:         c.StringSlice("dns"),
		PrivateIP:          c.String("bosh-private-ip"),
		PublicIP:           c.String("bosh-public-ip"),
		BoshReleaseSHA:     c.String("bosh-release-sha"),
		BoshReleaseURL:     c.String("bosh-release-url"),
		CPIReleaseSHA:      c.String("bosh-cpi-release-sha"),
		CPIReleaseURL:      c.String("bosh-cpi-release-url"),
		GOAgentSHA:         c.String("go-agent-release-sha"),
		GOAgentReleaseURL:  c.String("go-agent-release-url"),
		DirectorName:       c.String("director-name"),
		UAAReleaseSHA:      c.String("uaa-release-sha"),
		UAAReleaseURL:      c.String("uaa-release-url"),
		NtpServers:         c.StringSlice("ntp-server"),
		TrustedCerts:       c.String("trusted-certs"),
		NatsPassword:       c.String("nats-pwd"),
		PersistentDiskSize: c.Int("persistent-disk-size"),
		PrintManifest:      c.Bool("print-manifest"),
		GraphiteAddress:    c.String("hm-graphite-address"),
		GraphitePort:       c.Int("hm-graphite-port"),
		SyslogAddress:      c.String("syslog-address"),
		SyslogPort:         c.Int("syslog-port"),
		SyslogTransport:    c.String("syslog-transport"),
	}
	base.InitializePasswords()
	if base.IsUAA() {
		if err = base.InitializeCerts(); err != nil {
			return
		}
		if err = base.InitializeKeys(); err != nil {
			return
		}
	}
	return
}
