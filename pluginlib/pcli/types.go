package pcli

import (
	"strings"

	"github.com/codegangsta/cli"
)

type (
	Flag interface {
		ToCli() interface{}
	}
	StringFlag struct {
		Name   string
		Usage  string
		EnvVar string
		Value  string
	}
	StringSliceFlag struct {
		Name   string
		Usage  string
		EnvVar string
		Value  []string
	}
	BoolFlag struct {
		Name   string
		Usage  string
		EnvVar string
	}
	IntFlag struct {
		Name   string
		Usage  string
		EnvVar string
		Value  int
	}
	BoolTFlag struct {
		Name   string
		Usage  string
		EnvVar string
	}
)

func (s StringFlag) ToCli() interface{} {
	return cli.StringFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Value:  s.Value,
		Usage:  s.Usage,
	}
}

func (s StringSliceFlag) ToCli() interface{} {

	res := cli.StringSliceFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Value:  &cli.StringSlice{},
		Usage:  s.Usage,
	}
	res.Value.Set(strings.Join(s.Value, ","))
	return res
}

func (s BoolFlag) ToCli() interface{} {

	return cli.BoolFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Usage:  s.Usage,
	}
}

func (s IntFlag) ToCli() interface{} {

	return cli.IntFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Value:  s.Value,
		Usage:  s.Usage,
	}
}

func (s BoolTFlag) ToCli() interface{} {

	return cli.BoolTFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Usage:  s.Usage,
	}
}
