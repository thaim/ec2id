package main

import (
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

var (
	all     bool
	verbose bool
	version = ""
)

func main() {
	app := &cli.App{
		Name:  "ec2id",
		Usage: "get instance id",
		Flags: []cli.Flag{
			// &cli.BoolFlag{
			// 	Name: "help",
			// 	Destination: &help,
			// },
			&cli.BoolFlag{
				Name:        "all",
				Destination: &all,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Destination: &verbose,
			},
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.Args().Get(0)
			return Ec2id(name)
		},
		HideHelpCommand: true,
		Version: getVersion(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func getVersion() string {
	if version != "" {
		return version
	}
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	return i.Main.Version
}
