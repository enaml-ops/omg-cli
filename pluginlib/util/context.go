package pluginutil

import (
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
)

//NewContext - convenience method to construct a valid cli.Context within a
//plugin
func NewContext(args []string, myflags []cli.Flag) (context *cli.Context) {
	command := cli.Command{
		Flags: myflags,
	}
	app := cli.NewApp()
	app.Name = args[0]
	app.HideHelp = true
	app.Flags = myflags
	app.Action = func(c *cli.Context) error {
		c.Command = command
		context = c
		return nil
	}
	app.Run(args)
	return context
}

// LoadResourceFromContext loads data from a CLI flag.
//
// If the value of the specified flag starts with '@', the flag is interpreted
// as a filename and the contents of the file are returned.
//
// In all other cases, the flag value is returned directly.
func LoadResourceFromContext(c *cli.Context, flag string) (string, error) {
	value := c.String(flag)
	if len(value) > 0 && value[0] == '@' {
		b, err := ioutil.ReadFile(value[1:])
		if err != nil {
			return "", fmt.Errorf("couldn't read %s: %s\n", value[1:], err.Error())
		}
		value = string(b)
	}
	return value, nil
}
