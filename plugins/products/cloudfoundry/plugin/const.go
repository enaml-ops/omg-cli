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

	CFMysqlReleaseName         = "cf-mysql"
	CFLinuxReleaseName         = "cflinuxfs2-rootfs"
	EtcdReleaseName            = "etcd"
	PushAppsReleaseName        = "push-apps-manager-release"
	NotificationsReleaseName   = "notifications"
	NotificationsUIReleaseName = "notifications-ui"
	CFAutoscalingReleaseName   = "cf-autoscaling"
)

var (
	//DeploymentName -
	DeploymentName = "cf"
	//CFReleaseVersion -
	CFReleaseVersion = "235.5"
	//StemcellVersion -
	StemcellVersion = "3232.4"
	//DiegoReleaseVerion
	DiegoReleaseVersion = "0.1467.0"
	//CFMysqlReleaseVersion
	CFMysqlReleaseVersion = "25.2"

	GardenReleaseVersion          = "0.337.0"
	CFLinuxReleaseVersion         = "1.3.0"
	EtcdReleaseVersion            = "48"
	PushAppsReleaseVersion        = "621"
	NotificationsReleaseVersion   = "19"
	NotificationsUIReleaseVersion = "10"
	CFAutoscalingReleaseVersion   = "28"
)

var factories []InstanceGrouperFactory
