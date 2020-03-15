package main

import (
	"fmt"
	"log"
	"os"
	"space"

	"github.com/urfave/cli/v2"
)

func handleEnum(val string, enums []string) (value string, err error) {
	for _, enum := range enums {
		if val == enum {
			return val, err
		}
	}
	return "", fmt.Errorf("Invalid argument %v, possible values: %v", val, enums)
}

func handleEnvFlag(val string) (string, error) {
	return handleEnum(c.String("env"), []string{
		"dev",
		"live",
	})
}

func pushFolder(s Space, env, folder, prefix string) {
}

func pushAction(c *cli.Context) error {
	env, err := handleEnvFlag(c.String("env"))
	if err != nil {
		return err
	}

	s := space.New(service.SPACE_ENDPOINT)

	fp := c.Args().Get(0)
	prefix := c.String("prefix")
	if c.Bool("recursive") {
		return pushFolder(env, fp, prefix)
	}
	return pushFile(env, fp, prefix)
}

func main() {
	envFlag := cli.StringFlag{
		Name:     "env",
		Value:    "dev",
		Usage:    "Specify Space environment",
		Required: true,
	}

	pushCommand := cli.Command{
		Name:      "push",
		Usage:     "Upload file/folder to Space",
		ArgsUsage: "File or folder path to upload",
		Flags: []cli.Flag{
			&envFlag,
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "Upload a folder recursively",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "prefix",
				Aliases: []string{"p"},
				Usage:   "Object name's prefix.",
				Value:   "",
			},
		},
		Action: pushAction,
	}

	app := &cli.App{
		Name:  "space",
		Usage: "Work with Space and assets",
		Commands: []*cli.Command{
			&pushCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
