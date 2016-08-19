package boshinit

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/utils"
)

func BoshFlags(defaults *BoshBase) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "mode", Value: "basic", Usage: "what type of bosh director to install.  Options are basic or uaa"},
		cli.StringFlag{Name: "cidr", Value: defaults.NetworkCIDR, Usage: "the network cidr range for your bosh deployment"},
		cli.StringFlag{Name: "gateway", Value: defaults.NetworkGateway, Usage: "the gateway ip"},
		cli.StringSliceFlag{Name: "dns", Value: utils.ConvertToCLIStringSliceFlag(defaults.NetworkDNS), Usage: "the dns ip"},
		cli.StringFlag{Name: "bosh-private-ip", Value: defaults.PrivateIP, Usage: "the private ip for the bosh vm to be created"},
		cli.StringFlag{Name: "bosh-public-ip", Usage: "the public ip for the bosh vm to be created"},
		cli.StringFlag{Name: "bosh-release-sha", Value: defaults.BoshReleaseSHA, Usage: "sha1 of the bosh release being used (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-release-url", Value: defaults.BoshReleaseURL, Usage: "url to bosh release"},
		cli.StringFlag{Name: "bosh-cpi-release-sha", Value: defaults.CPIReleaseSHA, Usage: "sha1 of the cpi release being used (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-cpi-release-url", Value: defaults.CPIReleaseURL, Usage: "url to bosh cpi release"},
		cli.StringFlag{Name: "go-agent-release-sha", Value: defaults.GOAgentSHA, Usage: "sha1 of the go agent being use (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-release-url", Value: defaults.GOAgentReleaseURL, Usage: "url to stemcell release"},
		cli.StringFlag{Name: "director-name", Value: "enaml-bosh", Usage: "the name of your director"},
		cli.StringFlag{Name: "uaa-release-sha", Value: "899f1e10f27e82ac524f1158a513392bbfabf2a0", Usage: "sha1 of the uaa release being used (found on bosh.io)"},
		cli.StringFlag{Name: "uaa-release-url", Value: "https://bosh.io/d/github.com/cloudfoundry/uaa-release?v=12.2", Usage: "url to uaa release"},
		cli.StringSliceFlag{Name: "ntp-server", Value: utils.ConvertToCLIStringSliceFlag(defaults.NtpServers), Usage: "ntp server address"},
		cli.StringFlag{Name: "trusted-certs", Usage: "trusted ssl certs"},
		cli.StringFlag{Name: "nats-pwd", Usage: "password for nats"},
		cli.IntFlag{Name: "persistent-disk-size", Value: defaults.PersistentDiskSize, Usage: "size of persistent disk"},
		cli.BoolFlag{Name: "print-manifest", Usage: "if you would simply like to output a manifest the set this flag as true."},
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
		NetworkDNS:         utils.ClearDefaultStringSliceValue(c.StringSlice("dns")...),
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
		NtpServers:         utils.ClearDefaultStringSliceValue(c.StringSlice("ntp-server")...),
		TrustedCerts:       c.String("trusted-certs"),
		NatsPassword:       c.String("nats-pwd"),
		PersistentDiskSize: c.Int("persistent-disk-size"),
		PrintManifest:      c.Bool("print-manifest"),
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
