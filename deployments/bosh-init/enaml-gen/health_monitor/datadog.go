package health_monitor 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Datadog struct {

	/*ApplicationKey - Descr: Health Monitor Application Key for DataDog Default: <nil>
*/
	ApplicationKey interface{} `yaml:"application_key,omitempty"`

	/*PagerdutyServiceName - Descr: Service name to alert in PagerDuty upon HM events Default: <nil>
*/
	PagerdutyServiceName interface{} `yaml:"pagerduty_service_name,omitempty"`

	/*ApiKey - Descr: API Key for DataDog Default: <nil>
*/
	ApiKey interface{} `yaml:"api_key,omitempty"`

}