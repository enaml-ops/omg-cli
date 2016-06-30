package pluginutil

import "github.com/codegangsta/cli"

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
