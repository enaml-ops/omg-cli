package concourseplugin

import (
	"strings"

	"github.com/codegangsta/cli"
)

type flagBucket struct {
	Desc        string
	EnvVar      string
	StringSlice bool
}

const (
	defaultFileName              string = "concourse.yml"
	outputFileName               string = "OUTPUT_FILENAME"
	concoursePassword            string = "PASSWORD"
	concourseUsername            string = "USERNAME"
	concourseURL                 string = "URL"
	concourseWebInstances        string = "WEB_INSTANCES"
	concourseWebIPs              string = "WEB_IPS"
	boshDirectorUUID             string = "BOSH_DIRECTOR_UUID"
	boshStemcellAlias            string = "BOSH_STEMCELL_ALIAS"
	concourseNetworkName         string = "NETWORK_NAME"
	concourseNetworkRange        string = "NETWORK_RANGE"
	concourseNetworkGateway      string = "NETWORK_GATEWAY"
	concourseWebAZs              string = "WEB_AZS"
	concourseDatabaseAZs         string = "DATABASE_AZS"
	concourseWorkerAZs           string = "WORKER_AZS"
	boshCloudConfig              string = "BOSH_CLOUD_CONFIG"
	concourseDeploymentName      string = "BOSH_DEPLOYMENT_NAME"
	concourseWebVMType           string = "WEB_VM_TYPE"
	concourseDatabaseVMType      string = "DATABASE_VM_TYPE"
	concourseWorkerVMType        string = "WORKER_VM_TYPE"
	concourseDatabaseStorageType string = "DATABASE_STORAGE_TYPE"
	concoursePostgresqlDbPwd     string = "POSTGRESQL_DB_PWD"
	cloudConfigYml               string = "CLOUD_CONFIG_YML"
)

func getFlag(input string) (flag string) {
	flag = strings.ToLower(strings.Replace(input, "_", "-", -1))
	return
}

func generateFlags() (flags []cli.Flag) {
	var flagList = map[string]flagBucket{
		outputFileName: flagBucket{
			Desc:   "destination for output",
			EnvVar: outputFileName,
		},
		boshCloudConfig: flagBucket{
			Desc:   "true/false for generate cloudConfig compatible (default true)",
			EnvVar: boshCloudConfig,
		},
		boshDirectorUUID: flagBucket{
			Desc:   "bosh director uuid (bosh status --uuid)",
			EnvVar: boshDirectorUUID,
		},
		boshStemcellAlias: flagBucket{
			Desc:   "url for concourse ui",
			EnvVar: boshStemcellAlias,
		},
		concourseWebInstances: flagBucket{
			Desc:   "number of web instances (default 1)",
			EnvVar: concourseWebInstances,
		},
		concourseWebIPs: flagBucket{
			Desc:        "array of IPs reserved for concourse web-ui",
			EnvVar:      concourseWebIPs,
			StringSlice: true,
		},
		concourseURL: flagBucket{
			Desc:   "url for concourse ui",
			EnvVar: concourseURL,
		},
		concourseUsername: flagBucket{
			Desc:   "username for concourse ui",
			EnvVar: concourseUsername,
		},
		concoursePassword: flagBucket{
			Desc:   "password for concourse ui",
			EnvVar: concoursePassword,
		},
		concourseNetworkName: flagBucket{
			Desc:   "name of network to deploy concourse on",
			EnvVar: concourseNetworkName,
		},
		concourseNetworkRange: flagBucket{
			Desc:   "network range to deploy concourse on - only applies in non-cloud config",
			EnvVar: concourseNetworkRange,
		},
		concourseNetworkGateway: flagBucket{
			Desc:   "network gateway for concourse - only applies in non-cloud config",
			EnvVar: concourseNetworkGateway,
		},
		concourseWebAZs: flagBucket{
			Desc:        "array of AZs to deploy concourse web jobs to",
			EnvVar:      concourseWebAZs,
			StringSlice: true,
		},
		concourseDatabaseAZs: flagBucket{
			Desc:        "array of AZs to deploy concourse database jobs to",
			EnvVar:      concourseDatabaseAZs,
			StringSlice: true,
		},
		concourseWorkerAZs: flagBucket{
			Desc:        "array of AZs to deploy concourse worker jobs to",
			EnvVar:      concourseWorkerAZs,
			StringSlice: true,
		},
		concourseDeploymentName: flagBucket{
			Desc:   "bosh deployment name",
			EnvVar: concourseDeploymentName,
		},
		concoursePostgresqlDbPwd: flagBucket{
			Desc:   "password for postgresql job",
			EnvVar: concoursePostgresqlDbPwd,
		},
		concourseWebVMType: flagBucket{
			Desc:   "web vm type reference from cloudConfig",
			EnvVar: concourseWebVMType,
		},
		concourseDatabaseVMType: flagBucket{
			Desc:   "database vm type reference from cloudConfig",
			EnvVar: concourseDatabaseVMType,
		},
		concourseWorkerVMType: flagBucket{
			Desc:   "worker vm type reference from cloudConfig",
			EnvVar: concourseWorkerVMType,
		},
		concourseDatabaseStorageType: flagBucket{
			Desc:   "database storage type reference from cloudConfig",
			EnvVar: concourseDatabaseStorageType,
		},
		cloudConfigYml: flagBucket{
			Desc:   "location of cloud config yml",
			EnvVar: cloudConfigYml,
		},
	}
	for _, v := range flagList {
		if v.StringSlice {
			flags = append(flags, cli.StringSliceFlag{
				Name:   getFlag(v.EnvVar),
				Usage:  v.Desc,
				EnvVar: v.EnvVar,
			})
		} else {
			flags = append(flags, cli.StringFlag{
				Name:   getFlag(v.EnvVar),
				Value:  "",
				Usage:  v.Desc,
				EnvVar: v.EnvVar,
			})
		}
	}
	return
}
