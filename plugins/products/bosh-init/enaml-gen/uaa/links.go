package uaa 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Links struct {

	/*Passwd - Descr: URL for requesting password reset Default: /forgot_password
*/
	Passwd interface{} `yaml:"passwd,omitempty"`

	/*Signup - Descr: URL for requesting to signup/register for an account Default: /create_account
*/
	Signup interface{} `yaml:"signup,omitempty"`

}