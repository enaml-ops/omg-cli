package syslog_configurator 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type SyslogConfiguratorJob struct {

	/*SyslogAggregator - Descr: Transport to be used when forwarding logs (tcp|udp|relp). Default: udp
*/
	SyslogAggregator *SyslogAggregator `yaml:"syslog_aggregator,omitempty"`

}