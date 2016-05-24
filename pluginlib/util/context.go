package pluginutil

import (
	"sync"

	"github.com/codegangsta/cli"
)

//NewContext - convenience method to construct a valid cli.Context within a
//plugin
func NewContext(args []string, myflags []cli.Flag) (context *cli.Context) {
	var wg sync.WaitGroup
	app := cli.NewApp()
	app.Name = args[0]
	app.HideHelp = true
	app.Flags = myflags
	app.Action = func(c *cli.Context) error {
		defer wg.Done()
		context = c
		return nil
	}

	wg.Add(1)
	app.Run(args)
	wg.Wait()
	return context
}
