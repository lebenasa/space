package main

import (
	"fmt"
	"log"
	"os"

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

func main() {
	app := &cli.App{
		Name:  "space",
		Usage: "Upload assets to Space",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "Specify Space environment",
			},
		},
		Action: func(c *cli.Context) (err error) {
			env, err := handleEnum(c.String("env"), []string{
				"dev",
				"live",
			})
			if err != nil {
				return
			}
			fmt.Printf("Space environment: %v\n", env)
			return
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
