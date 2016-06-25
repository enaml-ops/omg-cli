package cloud_controller_worker 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type CloudControllerWorker struct {

	/*Hm9000 - Descr: URL of the hm9000 server Default: <nil>
*/
	Hm9000 *Hm9000 `yaml:"hm9000,omitempty"`

	/*Ccdb - Descr: The timeout for Sequel pooled connections Default: 10
*/
	Ccdb *Ccdb `yaml:"ccdb,omitempty"`

	/*Name - Descr: 'name' attribute in the /v2/info endpoint Default: 
*/
	Name interface{} `yaml:"name,omitempty"`

	/*SupportAddress - Descr: 'support' attribute in the /v2/info endpoint Default: 
*/
	SupportAddress interface{} `yaml:"support_address,omitempty"`

	/*Build - Descr: 'build' attribute in the /v2/info endpoint Default: 
*/
	Build interface{} `yaml:"build,omitempty"`

	/*LoggerEndpoint - Descr: Port for logger endpoint listed at /v2/info Default: 443
*/
	LoggerEndpoint *LoggerEndpoint `yaml:"logger_endpoint,omitempty"`

	/*Version - Descr: 'version' attribute in the /v2/info endpoint Default: 0
*/
	Version interface{} `yaml:"version,omitempty"`

	/*MetronEndpoint - Descr: The port used to emit messages to the Metron agent Default: 3457
*/
	MetronEndpoint *MetronEndpoint `yaml:"metron_endpoint,omitempty"`

	/*Nats - Descr: IP of each NATS cluster member. Default: <nil>
*/
	Nats *Nats `yaml:"nats,omitempty"`

	/*NfsServer - Descr: NFS server for droplets and apps (not used in an AWS deploy, use s3 instead) Default: <nil>
*/
	NfsServer *NfsServer `yaml:"nfs_server,omitempty"`

	/*Description - Descr: 'description' attribute in the /v2/info endpoint Default: 
*/
	Description interface{} `yaml:"description,omitempty"`

	/*AppDomains - Descr: Array of domains for user apps (example: 'user.app.space.foo', a user app called 'neat' will listen at 'http://neat.user.app.space.foo') Default: <nil>
*/
	AppDomains interface{} `yaml:"app_domains,omitempty"`

	/*RequestTimeoutInSeconds - Descr: Timeout for requests in seconds. Default: 900
*/
	RequestTimeoutInSeconds interface{} `yaml:"request_timeout_in_seconds,omitempty"`

	/*Domain - Descr: domain where cloud_controller will listen (api.domain) often the same as the system domain Default: <nil>
*/
	Domain interface{} `yaml:"domain,omitempty"`

	/*SystemDomainOrganization - Descr: The User Org that owns the system_domain, required if system_domain is defined Default: 
*/
	SystemDomainOrganization interface{} `yaml:"system_domain_organization,omitempty"`

	/*Ssl - Descr: specifies that the job is allowed to skip ssl cert verification Default: false
*/
	Ssl *Ssl `yaml:"ssl,omitempty"`

	/*SystemDomain - Descr: Domain reserved for CF operator, base URL where the login, uaa, and other non-user apps listen Default: <nil>
*/
	SystemDomain interface{} `yaml:"system_domain,omitempty"`

	/*DeaNext - Descr: Memory limit in mb for staging tasks Default: 1024
*/
	DeaNext *DeaNext `yaml:"dea_next,omitempty"`

	/*Cc - Descr: Specifies interval on which the CC will poll a service broker for asynchronous actions Default: 60
*/
	Cc *Cc `yaml:"cc,omitempty"`

	/*Login - Descr: whether use login as the authorization endpoint or not Default: true
*/
	Login *Login `yaml:"login,omitempty"`

	/*Uaa - Descr: Used for generating SSO clients for service brokers. Default: <nil>
*/
	Uaa *Uaa `yaml:"uaa,omitempty"`

}