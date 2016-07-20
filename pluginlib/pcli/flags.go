package pcli

import "strings"

type FlagType int

const (
	StringFlag FlagType = iota
	StringSliceFlag
	BoolFlag
	BoolTFlag // A boolean flag that is true by default.
	IntFlag
)

type Flag struct {
	Typ    FlagType
	Name   string
	Usage  string
	EnvVar string
	Value  string
}

// NewFlag creates a new flag with the specified data.
// The flag can be overridden with an environment variable
// of the same name.  For example, the flag 'my-cool-flag'
// can be overridden with the environment variable
// 'MY_COOL_FLAG'.
func NewFlag(t FlagType, name, usage, value string) Flag {
	return Flag{
		Typ:    t,
		Name:   name,
		Usage:  usage,
		EnvVar: flagToEnv(name),
		Value:  value,
	}
}

func flagToEnv(flag string) string {
	upper := strings.ToUpper(flag)
	return strings.Replace(upper, "-", "_", -1)
}
