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
	concourseDeploymentName      string = "BOSH_DEPLOYMENT_NAME"
	concourseWebVMType           string = "WEB_VM_TYPE"
	concourseDatabaseVMType      string = "DATABASE_VM_TYPE"
	concourseWorkerVMType        string = "WORKER_VM_TYPE"
	concourseDatabaseStorageType string = "DATABASE_STORAGE_TYPE"
	concoursePostgresqlDbPwd     string = "POSTGRESQL_DB_PWD"
	remoteStemcellURL            string = "REMOTE_STEMCELL_URL"
	remoteStemcellSHA            string = "REMOTE_STEMCELL_SHA"
)

func getFlag(input string) (flag string) {
	flag = strings.ToLower(strings.Replace(input, "_", "-", -1))
	return
}

func generateFlags() (flags []cli.Flag) {
	var flagList = map[string]flagBucket{
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
		remoteStemcellURL: flagBucket{
			Desc:   "url to the remote stemcell you wish to use.",
			EnvVar: remoteStemcellURL,
		},
		remoteStemcellSHA: flagBucket{
			Desc:   "sha1 of the remote stemcell.",
			EnvVar: remoteStemcellSHA,
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
