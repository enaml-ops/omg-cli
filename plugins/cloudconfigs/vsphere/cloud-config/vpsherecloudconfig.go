package cloudconfig

import (
	"errors"
	"fmt"
)

func (c *VSphereCloudConfig) validate() error {
	if len(c.AZs) == 0 {
		return errors.New("One or more availability zones are required")
	}
	for i, az := range c.AZs {
		if len(az.Name) == 0 {
			return fmt.Errorf("AZ %d must have a name", i)
		}
		if len(az.Cluster.Name) == 0 {
			return fmt.Errorf("AZ '%s' must have a vSphere cluster name", az.Name)
		}
		if len(az.Cluster.ResourcePool) == 0 {
			return fmt.Errorf("AZ '%s' must have a vSphere resource pool name", az.Name)
		}
		if len(az.Network.Name) == 0 {
			return fmt.Errorf("AZ '%s' must have a vSphere network name", az.Name)
		}
		if len(az.Network.Range) == 0 {
			return fmt.Errorf("AZ '%s' must have a CIDR formatted IP range", az.Name)
		}
		if len(az.Network.Gateway) == 0 {
			return fmt.Errorf("AZ '%s' must have a gateway IP", az.Name)
		}
		if len(az.Network.DNS) == 0 {
			return fmt.Errorf("AZ '%s' must have one or more DNS servers", az.Name)
		}
	}
	return nil
}

// VSphereCloudConfig contains all the necessary information to build a vSphere
// cloud config manifest
type VSphereCloudConfig struct {
	AZs []VSphereAZ
}

// VSphereAZ holds the cloud config availability zone details for vSphere
type VSphereAZ struct {
	Name    string
	Cluster VSphereCluster
	Network VSphereNetwork
}

// VSphereCluster ties a vSphere cluster to a vSphere resource pool
type VSphereCluster struct {
	Name         string
	ResourcePool string
}

// VSphereNetwork is a vSphere subnet
type VSphereNetwork struct {
	Name    string
	Range   string
	Gateway string
	DNS     []string
}
