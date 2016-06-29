package cloudfoundry

const (
	CFReleaseName = "cf"
	StemcellName  = "ubuntu-trusty"
	StemcellAlias = "trusty"
)

var (
	DeploymentName   = "cf"
	CFReleaseVersion = "235.5"
	StemcellVersion  = "3232.4"
)

var factories []InstanceGrouperFactory
