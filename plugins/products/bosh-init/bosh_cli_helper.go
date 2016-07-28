package boshinit

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func BoshFlags(defaults BoshDefaults) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "mode", Value: "basic", Usage: "what type of bosh director to install.  Options are basic or uaa"},
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
		cli.StringFlag{Name: "director-name", Value: "enaml-bosh", Usage: "the name of your director"},
		cli.StringFlag{Name: "uaa-release-ver", Value: "12.2", Usage: "the bosh uaa version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "uaa-release-sha", Value: "899f1e10f27e82ac524f1158a513392bbfabf2a0", Usage: "sha1 of the uaa release being used (found on bosh.io)"},
		cli.StringFlag{Name: "cpi-name", Value: defaults.CPIName, Usage: ""},
		cli.StringSliceFlag{Name: "ntp-server", Value: defaults.NtpServers, Usage: ""},
		cli.BoolFlag{Name: "print-manifest", Usage: "if you would simply like to output a manifest the set this flag as true."},
	}
}

func checkRequired(names []string, c *cli.Context) {
	var invalidNames []string
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
		"uaa-release-sha",
		"cpi-name", "ntp-server",
	}, c)

	base = &BoshBase{
		Mode:               c.String("mode"),
		NetworkCIDR:        c.String("cidr"),
		NetworkGateway:     c.String("gateway"),
		NetworkDNS:         c.StringSlice("dns"),
		PrivateIP:          c.String("bosh-private-ip"),
		PublicIP:           c.String("bosh-public-ip"),
		BoshReleaseVersion: c.String("bosh-release-ver"),
		BoshReleaseSHA:     c.String("bosh-release-sha"),
		CPIReleaseVersion:  c.String("bosh-cpi-release-ver"),
		CPIReleaseSHA:      c.String("bosh-cpi-release-sha"),
		GOAgentVersion:     c.String("go-agent-ver"),
		GOAgentSHA:         c.String("go-agent-sha"),
		DirectorName:       c.String("director-name"),
		UAAReleaseVersion:  c.String("uaa-release-ver"),
		UAAReleaseSHA:      c.String("uaa-release-sha"),
		CPIName:            c.String("cpi-name"),
		NtpServers:         c.StringSlice("ntp-server"),
	}
	base.InitializePasswords()
	fmt.Println("**********************************")
	if base.IsUAA() {
		if err = base.InitializeCerts(); err != nil {
			return
		}
		if err = base.InitializeKeys(); err != nil {
			return
		}
		fmt.Println("Director CA Certificate")
		fmt.Println(base.CACert)
	}
	fmt.Println("Director PWD:", base.DirectorPassword)
	fmt.Println("**********************************")
	return
}
