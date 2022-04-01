package main

import (
	"os"

	"git.querycap.com/musenwill/memalloc/alloc"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.ErrWriter = os.Stderr
	app.EnableBashCompletion = true
	app.Name = "memalloc"
	app.Usage = "Memory alloc experiment"
	app.Flags = []cli.Flag{maxLimitFlag, minSizeFlag, maxSizeFlag, spreadFlag, reMinSizeFlag, reMaxSizeFlag, reSpreadFlag, printFlag}
	app.Action = action
	app.RunAndExitOnError()
}

func action(c *cli.Context) error {
	a := alloc.NewAlloc(alloc.Config{
		MaxLimit:  c.GlobalInt64(maxLimitFlag.Name) * 1024 * 1024,
		MinSize:   c.GlobalInt64(minSizeFlag.Name),
		MaxSize:   c.GlobalInt64(maxSizeFlag.Name),
		Spread:    c.GlobalFloat64(spreadFlag.Name),
		ReMinSize: c.GlobalInt64(reMinSizeFlag.Name),
		ReMaxSize: c.GlobalInt64(reMaxSizeFlag.Name),
		ReSpread:  c.GlobalFloat64(reSpreadFlag.Name),
		Print:     c.GlobalBool(printFlag.Name),
	})
	a.Run()
	return nil
}
