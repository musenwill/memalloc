package main

import "github.com/urfave/cli"

var maxLimitFlag = cli.Int64Flag{
	Name:     "max-limit",
	Usage:    "max memory limit in MB",
	Required: false,
	Value:    2 * 1024,
}

var minSizeFlag = cli.Int64Flag{
	Name:     "min",
	Usage:    "min size in bytes",
	Required: true,
}

var maxSizeFlag = cli.Int64Flag{
	Name:     "max",
	Usage:    "max size in bytes",
	Required: true,
}

var spreadFlag = cli.Float64Flag{
	Name:     "spread",
	Usage:    "data distribution",
	Required: false,
	Value:    0.5,
}

var reMinSizeFlag = cli.Int64Flag{
	Name:     "rmin",
	Usage:    "min size in bytes",
	Required: true,
}

var reMaxSizeFlag = cli.Int64Flag{
	Name:     "rmax",
	Usage:    "max size in bytes",
	Required: true,
}

var reSpreadFlag = cli.Float64Flag{
	Name:     "rspread",
	Usage:    "data distribution",
	Required: false,
	Value:    0.5,
}

var printFlag = cli.BoolFlag{
	Name:     "print",
	Usage:    "print csv",
	Required: false,
}

var halfFlag = cli.BoolFlag{
	Name:     "half",
	Usage:    "release half memory",
	Required: false,
}
