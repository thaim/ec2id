package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

var (
	all     bool
	verbose bool
	version = ""
	revision = ""
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
			client, err := NewAwsClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "initialized failed:%v\n", err)
				os.Exit(1)
			}
			ids, err := Ec2id(name, client)
			if err == nil && ids != nil {
				printIds(os.Stdout, ids, all)
			}
			return err
		},
		HideHelpCommand: true,
		Version:         versionFormatter(getVersion(), getRevision()),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func versionFormatter(version string, revision string) string {
	if version == "" {
		version = "devel"
	}

	if revision == "" {
		return version
	}
	return fmt.Sprintf("%s (rev: %s)", version, revision)
}

func getVersion() string {
	if version != "" {
		return version
	}
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	return i.Main.Version
}

func getRevision() string {
	if revision != "" {
		return revision
	}
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	for _, s := range i.Settings {
		if s.Key == "vcs.revision" {
			return s.Value
		}
	}

	return ""
}

func printIds(out io.Writer, ids []string, all bool) {
	for _, id := range ids {
		fmt.Fprintln(out, id)

		if !all {
			break
		}
	}
}
