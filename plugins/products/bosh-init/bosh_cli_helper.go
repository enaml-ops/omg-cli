package boshinit

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func BoshFlags(defaults BoshDefaults) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "cidr", Value: defaults.CIDR, Usage: "the network cidr range for your bosh deployment"},
		cli.StringFlag{Name: "gateway", Value: defaults.Gateway, Usage: "the gateway ip"},
		cli.StringSliceFlag{Name: "dns", Value: defaults.DNS, Usage: "the dns ip"},
		cli.StringFlag{Name: "bosh-private-ip", Value: defaults.PrivateIP, Usage: "the private ip for the bosh vm to be created"},
		cli.StringFlag{Name: "bosh-public-ip", Usage: "the public ip for the bosh vm to be created"},
		cli.StringFlag{Name: "bosh-release-ver", Value: defaults.BoshReleaseVersion, Usage: "the version of the bosh release you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-release-sha", Value: defaults.BoshReleaseSHA, Usage: "sha1 of the bosh release being used (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-cpi-release-ver", Value: defaults.CPIReleaseVersion, Usage: "the bosh cpi version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-cpi-release-sha", Value: defaults.CPIReleaseSHA, Usage: "sha1 of the cpi release being used (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-ver", Value: defaults.GOAgentVersion, Usage: "the go agent version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-sha", Value: defaults.GOAgentSHA, Usage: "sha1 of the go agent being use (found on bosh.io)"},
		cli.StringFlag{Name: "director-name", Value: "my-bosh", Usage: "the name of your director"},
		cli.StringFlag{Name: "uaa-release-ver", Value: "12.2", Usage: "the bosh uaa version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "uaa-release-sha", Value: "899f1e10f27e82ac524f1158a513392bbfabf2a0", Usage: "sha1 of the uaa release being used (found on bosh.io)"},
		cli.StringFlag{Name: "uaa-jwt-signing-key", Usage: ""},
		cli.StringFlag{Name: "uaa-jwt-verification-key", Usage: ""},
		cli.StringFlag{Name: "uaa-public-key", Usage: ""},
		cli.StringFlag{Name: "director-password", Usage: ""},
		cli.StringFlag{Name: "agent-password", Usage: ""},
		cli.StringFlag{Name: "db-password", Usage: ""},
		cli.StringFlag{Name: "cpi-name", Value: defaults.CPIName, Usage: ""},
		cli.StringSliceFlag{Name: "ntp-server", Value: defaults.NtpServers, Usage: ""},
		cli.StringFlag{Name: "nats-password", Usage: ""},
		cli.StringFlag{Name: "mbus-password", Usage: ""},
		cli.StringFlag{Name: "ssl-cert", Usage: ""},
		cli.StringFlag{Name: "ssl-key", Usage: ""},
		cli.StringFlag{Name: "login-secret", Usage: ""},
		cli.StringFlag{Name: "registry-password", Usage: ""},
		cli.StringFlag{Name: "health-monitor-secret", Usage: ""},
		cli.StringFlag{Name: "health-monitor-ca-cert", Usage: ""},
		cli.BoolFlag{Name: "print-manifest", Usage: "if you would simply like to output a manifest the set this flag as true."},
	}
}

func checkRequired(names []string, c *cli.Context) {
	invalidNames := []string{""}
	for _, name := range names {
		if c.String(name) == "" {
			invalidNames = append(invalidNames, name)
		}
	}
	if len(invalidNames) > 0 {
		fmt.Println("Sorry you need to provide", invalidNames, "flags to continue")
		os.Exit(1)
	}
}

func NewBoshBase(c *cli.Context) (base *BoshBase, err error) {

	checkRequired([]string{
		"cidr", "gateway", "dns", "bosh-private-ip",
		"bosh-release-ver", "bosh-release-sha", "bosh-cpi-release-ver", "bosh-cpi-release-sha",
		"go-agent-ver", "go-agent-sha", "director-name", "uaa-release-ver",
		"uaa-release-sha", "uaa-jwt-signing-key", "uaa-jwt-verification-key", "uaa-public-key",
		"director-password", "db-password", "cpi-name", "ntp-server",
		"nats-password", "mbus-password", "ssl-cert", "ssl-key",
		"login-secret", "registry-password", "health-monitor-secret", "health-monitor-ca-cert",
	}, c)

	base = &BoshBase{
		NetworkCIDR:         c.String("cidr"),
		NetworkGateway:      c.String("gateway"),
		NetworkDNS:          c.StringSlice("dns"),
		PrivateIP:           c.String("bosh-private-ip"),
		PublicIP:            c.String("bosh-public-ip"),
		BoshReleaseVersion:  c.String("bosh-release-ver"),
		BoshReleaseSHA:      c.String("bosh-release-sha"),
		CPIReleaseVersion:   c.String("bosh-cpi-release-ver"),
		CPIReleaseSHA:       c.String("bosh-cpi-release-sha"),
		GOAgentVersion:      c.String("go-agent-ver"),
		GOAgentSHA:          c.String("go-agent-sha"),
		DirectorName:        c.String("director-name"),
		UAAReleaseVersion:   c.String("uaa-release-ver"),
		UAAReleaseSHA:       c.String("uaa-release-sha"),
		SigningKey:          c.String("uaa-jwt-signing-key"),
		VerificationKey:     c.String("uaa-jwt-verification-key"),
		UAAPublicKey:        c.String("uaa-public-key"),
		DirectorPassword:    c.String("director-password"),
		AgentPassword:       c.String("agent-password"),
		DBPassword:          c.String("db-password"),
		CPIName:             c.String("cpi-name"),
		NtpServers:          c.StringSlice("ntp-server"),
		NatsPassword:        c.String("nats-password"),
		MBusPassword:        c.String("mbus-password"),
		SSLCert:             c.String("ssl-cert"),
		SSLKey:              c.String("ssl-key"),
		LoginSecret:         c.String("login-secret"),
		RegistryPassword:    c.String("registry-password"),
		HealthMonitorSecret: c.String("health-monitor-secret"),
		CACert:              c.String("health-monitor-ca-cert"),
	}
	if base.PublicIP == "" {
		base.PublicIP = base.PrivateIP
	}
	return
}
