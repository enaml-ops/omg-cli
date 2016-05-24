package concourse

//Garden -
type Garden struct {
	Garden interface{} `yaml:"garden,omitempty"`
}

//DBName -
type DBName struct {
	Name     string `yaml:"name,omitempty"`
	Role     string `yaml:"role,omitempty"`
	Password string `yaml:"password,omitempty"`
}
