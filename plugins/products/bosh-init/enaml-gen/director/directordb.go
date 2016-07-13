package director 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type DirectorDb struct {

	/*User - Descr: Username used for the director database Default: bosh
*/
	User interface{} `yaml:"user,omitempty"`

	/*Host - Descr: Address of the director database, for example, in the case of AWS RDS:
rds-instance-name.coqxxxxxxxxx.us-east-1.rds.amazonaws.com
 Default: 127.0.0.1
*/
	Host interface{} `yaml:"host,omitempty"`

	/*ConnectionOptions - Descr: Additional options for the database Default: map[pool_timeout:10 max_connections:32]
*/
	ConnectionOptions interface{} `yaml:"connection_options,omitempty"`

	/*Password - Descr: Password used for the director database Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*Database - Descr: Name of the director database Default: bosh
*/
	Database interface{} `yaml:"database,omitempty"`

	/*Port - Descr: Port of the director database (e.g, msyql2 adapter would generally use 3306) Default: 5432
*/
	Port interface{} `yaml:"port,omitempty"`

	/*Adapter - Descr: The type of database used (mysql2|postgres|sqlite) Default: postgres
*/
	Adapter interface{} `yaml:"adapter,omitempty"`

}