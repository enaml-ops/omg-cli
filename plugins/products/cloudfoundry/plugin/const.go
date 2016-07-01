package cloudfoundry

const (
	//CFReleaseName -
	CFReleaseName = "cf"
	//StemcellName -
	StemcellName = "ubuntu-trusty"
	//StemcellAlias -
	StemcellAlias = "trusty"

	//CFLinuxFSReleaseName -
	CFLinuxFSReleaseName = "cflinuxfs2-rootfs"

	//GardenReleaseName
	GardenReleaseName = "garden-linux"

	//DiegoReleaseName
	DiegoReleaseName = "diego"
)

var (
	//DeploymentName -
	DeploymentName = "cf"
	//CFReleaseVersion -
	CFReleaseVersion = "235.5"
	//StemcellVersion -
	StemcellVersion = "3232.4"
)

var factories []InstanceGrouperFactory
