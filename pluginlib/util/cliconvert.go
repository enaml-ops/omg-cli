package pluginutil

import (
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-cli/pluginlib/pcli"
)

func ToCliFlagArray(fs []pcli.Flag) []cli.Flag {
	result := make([]cli.Flag, 0, len(fs))
	for i := range fs {
		result = append(result, toCLI(&fs[i]))
	}
	return result
}

func toCLI(f *pcli.Flag) cli.Flag {
	switch f.Typ {
	case pcli.StringFlag:
		return cli.StringFlag{
			Name:   f.Name,
			EnvVar: f.EnvVar,
			Value:  f.Value,
			Usage:  f.Usage,
		}
	case pcli.StringSliceFlag:
		ss := cli.StringSliceFlag{
			Name:   f.Name,
			EnvVar: f.EnvVar,
			Value:  &cli.StringSlice{},
			Usage:  f.Usage,
		}
		if f.Value != "" {
			ss.Value.Set(f.Value)
		}
		return ss
	case pcli.IntFlag:
		flag := cli.IntFlag{
			Name:   f.Name,
			EnvVar: f.EnvVar,
			Usage:  f.Usage,
		}
		if f.Value != "" {
			i, err := strconv.Atoi(f.Value)
			if err != nil {
				panic("Invalid int flag: " + f.Value)
			}
			flag.Value = i
		}
		return flag
	case pcli.BoolFlag:
		return cli.BoolFlag{
			Name:   f.Name,
			EnvVar: f.EnvVar,
			Usage:  f.Usage,
		}
	case pcli.BoolTFlag:
		return cli.BoolTFlag{
			Name:   f.Name,
			EnvVar: f.EnvVar,
			Usage:  f.Usage,
		}
	default:
		panic("Unknown flag type: " + strconv.Itoa(int(f.Typ)))
	}
}
