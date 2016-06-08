package concourseplugin

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/concourse"
)

func NewDeploymentManifest(c *cli.Context, cloudConfig []byte) enaml.DeploymentManifest {
	var deployment = concourse.NewDeployment()

	if c.IsSet(getFlag(concourseDeploymentName)) {
		deployment.DeploymentName = c.String(getFlag(concourseDeploymentName))

	} else {
		deployment.DeploymentName = "concourse"
	}

	if c.IsSet(getFlag(concoursePostgresqlDbPwd)) {
		deployment.PostgresPassword = c.String(getFlag(concoursePostgresqlDbPwd))

	} else {
		deployment.PostgresPassword = "dummy-postgres-password"
	}
	deployment.StemcellURL = c.String(getFlag(remoteStemcellURL))
	deployment.StemcellSHA = c.String(getFlag(remoteStemcellSHA))
	deployment.ConcoursePassword = c.String(getFlag(concoursePassword))
	deployment.ConcourseUserName = c.String(getFlag(concourseUsername))
	deployment.ConcourseURL = c.String(getFlag(concourseURL))
	deployment.DirectorUUID = c.String(getFlag(boshDirectorUUID))
	deployment.StemcellAlias = c.String(getFlag(boshStemcellAlias))
	deployment.NetworkName = c.String(getFlag(concourseNetworkName))

	if c.IsSet(getFlag(boshCloudConfig)) {
		deployment.CloudConfig = c.Bool(getFlag(boshCloudConfig))

	} else {
		deployment.CloudConfig = true
	}

	if c.IsSet(getFlag(concourseWebInstances)) {
		deployment.WebInstances = c.Int(getFlag(concourseWebInstances))

	} else {
		deployment.WebInstances = 1

	}

	if c.IsSet(getFlag(concourseWebIPs)) {
		deployment.WebIPs = c.StringSlice(getFlag(concourseWebIPs))
	}
	deployment.NetworkRange = c.String(getFlag(concourseNetworkRange))
	deployment.NetworkGateway = c.String(getFlag(concourseNetworkGateway))

	if c.IsSet(getFlag(concourseWebAZs)) {
		deployment.WebAZs = c.StringSlice(getFlag(concourseWebAZs))
	}

	if c.IsSet(getFlag(concourseDatabaseAZs)) {
		deployment.DatabaseAZs = c.StringSlice(getFlag(concourseDatabaseAZs))
	}

	if c.IsSet(getFlag(concourseWorkerAZs)) {
		deployment.WorkerAZs = c.StringSlice(getFlag(concourseWorkerAZs))
	}
	deployment.WebVMType = c.String(getFlag(concourseWebVMType))
	deployment.WorkerVMType = c.String(getFlag(concourseWorkerVMType))
	deployment.DatabaseVMType = c.String(getFlag(concourseDatabaseVMType))
	deployment.DatabaseStorageType = c.String(getFlag(concourseDatabaseStorageType))
	deployment.CloudConfigYml = c.String(getFlag(cloudConfigYml))

	if err := deployment.Initialize(cloudConfig); err != nil {
		panic(err.Error())
	}
	return deployment.GetDeployment()
}
