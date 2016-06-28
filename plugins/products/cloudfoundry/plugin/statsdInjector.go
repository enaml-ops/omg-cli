package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
)

//NewStatsdInjector -
func NewStatsdInjector(c *cli.Context) (statsdInjector *StatsdInjector, err error) {
	statsdInjector = &StatsdInjector{}
	return
}

//CreateJob -
func (s *StatsdInjector) CreateJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "statsd-injector",
		Release:    "cf",
		Properties: make(map[interface{}]interface{}),
	}
}
